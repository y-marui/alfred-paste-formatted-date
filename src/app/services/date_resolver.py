"""date_resolver - parse a query string into a target datetime.

Supported query forms
---------------------
Relative offsets  (sign + number + optional unit):
    -2   -2d    two days ago
    +1w         one week later
    -3m         three months ago
    +1y         one year later
    Unit defaults to 'd' when omitted.

Direct dates  (year/month/day  or  year-month-day):
    2026/7/9    26/7/9    2026-7-9
    Two-digit years are interpreted as 2000+YY.

Combined form (date spec + format filter, space-separated):
    -2d ISO        two days ago, ISO formats only
    2026/7/9 unix  specific date, unix timestamp only
    +1w YYYY       one week later, formats containing "YYYY"

Anything else is treated as a format filter; today is used as the target.

Returns
-------
tuple[datetime.datetime, str]
    (target_datetime, remaining_filter_query)
"""

from __future__ import annotations

import calendar
import datetime
import re


def resolve_date(query: str) -> tuple[datetime.datetime, str]:
    """Return (target_datetime, remaining_filter_query) for the given query.

    The first whitespace-separated token is tested as a date spec.
    If it matches, the remainder of the query is returned as the filter.
    Otherwise the whole query is used as a format filter against today.
    """
    q = query.strip()
    if not q:
        return datetime.datetime.now(), ""

    first, _, rest = q.partition(" ")

    dt = _try_relative(first)
    if dt is not None:
        return dt, rest.strip()

    dt = _try_direct_date(first)
    if dt is not None:
        return dt, rest.strip()

    return datetime.datetime.now(), q


# ---------------------------------------------------------------------------
# Internal helpers
# ---------------------------------------------------------------------------

_RELATIVE_RE = re.compile(r"^([+-])(\d+)([dwmy]?)$", re.IGNORECASE)


def _try_relative(q: str) -> datetime.datetime | None:
    m = _RELATIVE_RE.match(q)
    if not m:
        return None

    sign = 1 if m.group(1) == "+" else -1
    amount = int(m.group(2)) * sign
    unit = m.group(3).lower() or "d"

    today = datetime.date.today()
    if unit == "d":
        target = today + datetime.timedelta(days=amount)
    elif unit == "w":
        target = today + datetime.timedelta(weeks=amount)
    elif unit == "m":
        target = _add_months(today, amount)
    else:  # y
        target = _add_months(today, amount * 12)

    return datetime.datetime.combine(target, datetime.datetime.now().time().replace(microsecond=0))


_DATE_RE = re.compile(r"^(\d{2,4})[/-](\d{1,2})[/-](\d{1,2})$")


def _try_direct_date(q: str) -> datetime.datetime | None:
    m = _DATE_RE.match(q)
    if not m:
        return None

    year, month, day = int(m.group(1)), int(m.group(2)), int(m.group(3))
    if year < 100:
        year += 2000

    try:
        return datetime.datetime.combine(datetime.date(year, month, day), datetime.time.min)
    except ValueError:
        return None


def _add_months(date: datetime.date, months: int) -> datetime.date:
    """Add (or subtract) a number of months, clamping the day to the last valid day."""
    total_months = date.month - 1 + months
    year = date.year + total_months // 12
    month = total_months % 12 + 1
    day = min(date.day, calendar.monthrange(year, month)[1])
    return datetime.date(year, month, day)
