"""Tests for date_resolver."""

from __future__ import annotations

import datetime

import pytest

from app.services.date_resolver import _add_months, resolve_date


class TestResolveDate:
    def test_empty_returns_today(self):
        dt, remaining = resolve_date("")
        assert dt.date() == datetime.date.today()
        assert remaining == ""

    def test_whitespace_returns_today(self):
        dt, remaining = resolve_date("   ")
        assert dt.date() == datetime.date.today()
        assert remaining == ""

    # --- relative: days ---

    def test_minus_days_no_unit(self):
        dt, remaining = resolve_date("-2")
        assert dt.date() == datetime.date.today() - datetime.timedelta(days=2)
        assert remaining == ""

    def test_plus_days_no_unit(self):
        dt, remaining = resolve_date("+3")
        assert dt.date() == datetime.date.today() + datetime.timedelta(days=3)
        assert remaining == ""

    def test_minus_days_with_unit(self):
        dt, remaining = resolve_date("-2d")
        assert dt.date() == datetime.date.today() - datetime.timedelta(days=2)
        assert remaining == ""

    def test_plus_days_with_unit(self):
        dt, remaining = resolve_date("+1d")
        assert dt.date() == datetime.date.today() + datetime.timedelta(days=1)
        assert remaining == ""

    # --- relative: weeks ---

    def test_plus_weeks(self):
        dt, remaining = resolve_date("+1w")
        assert dt.date() == datetime.date.today() + datetime.timedelta(weeks=1)
        assert remaining == ""

    def test_minus_weeks(self):
        dt, remaining = resolve_date("-2w")
        assert dt.date() == datetime.date.today() - datetime.timedelta(weeks=2)
        assert remaining == ""

    # --- relative: months ---

    def test_minus_months(self):
        dt, remaining = resolve_date("-3m")
        assert dt.date() < datetime.date.today()
        assert remaining == ""

    def test_plus_months(self):
        dt, remaining = resolve_date("+1m")
        assert dt.date() > datetime.date.today()
        assert remaining == ""

    # --- relative: years ---

    def test_plus_years(self):
        dt, remaining = resolve_date("+1y")
        assert dt.date() > datetime.date.today()
        assert remaining == ""

    def test_minus_years(self):
        dt, remaining = resolve_date("-1y")
        assert dt.date() < datetime.date.today()
        assert remaining == ""

    # --- direct date ---

    def test_direct_date_4digit_slash(self):
        dt, remaining = resolve_date("2026/7/9")
        assert dt.date() == datetime.date(2026, 7, 9)
        assert remaining == ""

    def test_direct_date_2digit_slash(self):
        dt, remaining = resolve_date("26/7/9")
        assert dt.date() == datetime.date(2026, 7, 9)
        assert remaining == ""

    def test_direct_date_4digit_hyphen(self):
        dt, remaining = resolve_date("2026-7-9")
        assert dt.date() == datetime.date(2026, 7, 9)
        assert remaining == ""

    def test_direct_date_zero_padded(self):
        dt, remaining = resolve_date("2026/07/09")
        assert dt.date() == datetime.date(2026, 7, 9)
        assert remaining == ""

    def test_direct_date_midnight(self):
        dt, _ = resolve_date("2026/7/9")
        assert dt.time() == datetime.time.min

    def test_invalid_date_passthrough(self):
        """Invalid calendar date falls through to format filter."""
        dt, remaining = resolve_date("2026/13/1")
        assert dt.date() == datetime.date.today()
        assert remaining == "2026/13/1"

    # --- format filter passthrough ---

    def test_format_filter_passthrough(self):
        dt, remaining = resolve_date("ISO")
        assert dt.date() == datetime.date.today()
        assert remaining == "ISO"

    def test_format_filter_preserves_case(self):
        _, remaining = resolve_date("YYYY")
        assert remaining == "YYYY"

    # --- combined: date spec + format filter ---

    def test_relative_with_filter(self):
        dt, remaining = resolve_date("-2d ISO")
        assert dt.date() == datetime.date.today() - datetime.timedelta(days=2)
        assert remaining == "ISO"

    def test_direct_date_with_filter(self):
        dt, remaining = resolve_date("2026/7/9 unix")
        assert dt.date() == datetime.date(2026, 7, 9)
        assert remaining == "unix"

    def test_weeks_with_filter(self):
        dt, remaining = resolve_date("+1w YYYY")
        assert dt.date() == datetime.date.today() + datetime.timedelta(weeks=1)
        assert remaining == "YYYY"

    def test_filter_only_is_not_confused_with_date_spec(self):
        """A plain filter string must not be consumed as a date spec."""
        dt, remaining = resolve_date("YYYYMMDD")
        assert dt.date() == datetime.date.today()
        assert remaining == "YYYYMMDD"


class TestAddMonths:
    @pytest.mark.parametrize(
        "base, months, expected",
        [
            (datetime.date(2026, 4, 14), 3, datetime.date(2026, 7, 14)),
            (datetime.date(2026, 4, 14), -3, datetime.date(2026, 1, 14)),
            (datetime.date(2026, 1, 31), 1, datetime.date(2026, 2, 28)),  # clamp to Feb 28
            (datetime.date(2024, 1, 31), 1, datetime.date(2024, 2, 29)),  # leap year
            (datetime.date(2026, 12, 31), 2, datetime.date(2027, 2, 28)),
        ],
    )
    def test_add_months(self, base: datetime.date, months: int, expected: datetime.date) -> None:
        assert _add_months(base, months) == expected
