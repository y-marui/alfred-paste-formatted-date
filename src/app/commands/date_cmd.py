"""date command - list date formats for selection and paste.

Usage in Alfred:
    date                  — today in all formats
    date -2  / -2d        — 2 days ago
    date +1w              — 1 week later
    date -3m              — 3 months ago
    date +1y              — 1 year later
    date 2026/7/9         — specific date (4- or 2-digit year, / or - separator)
    date <filter>         — filter by format name or value (e.g. "ISO", "YYYY", "unix")
"""

from __future__ import annotations

import datetime
from typing import NamedTuple

from alfred.logger import get_logger
from alfred.response import item, output
from app.services.date_resolver import resolve_date

log = get_logger(__name__)


class DateFormat(NamedTuple):
    uid: str
    label: str
    strftime: str | None


def _unix_timestamp(dt: datetime.datetime) -> str:
    return str(int(dt.timestamp()))


_FORMATS: list[DateFormat] = [
    DateFormat("yyyymmdd", "YYYYMMDD", "%Y%m%d"),
    DateFormat("yymmdd", "YYMMDD", "%y%m%d"),
    DateFormat("iso-date", "YYYY-MM-DD", "%Y-%m-%d"),
    DateFormat("slash-ymd", "YYYY/MM/DD", "%Y/%m/%d"),
    DateFormat("slash-mdy", "MM/DD/YYYY", "%m/%d/%Y"),
    DateFormat("slash-dmy", "DD/MM/YYYY", "%d/%m/%Y"),
    DateFormat("abbr-month", "MMM DD, YYYY", "%b %d, %Y"),
    DateFormat("full-month", "MMMM DD, YYYY", "%B %d, %Y"),
    DateFormat("iso-datetime", "YYYY-MM-DDThh:mm:ss", "%Y-%m-%dT%H:%M:%S"),
    DateFormat("unix", "Unix timestamp", None),
]


def _formatted(fmt: DateFormat, dt: datetime.datetime) -> str:
    if fmt.uid == "unix":
        return _unix_timestamp(dt)
    assert fmt.strftime is not None
    return dt.strftime(fmt.strftime)


def handle(args: str) -> None:
    """Resolve the target date from args, then show matching date formats."""
    log.debug("date command: args=%r", args)

    target_dt, filter_query = resolve_date(args)
    query = filter_query.strip().lower()

    items = []
    for fmt in _FORMATS:
        value = _formatted(fmt, target_dt)
        if query and not (query in fmt.label.lower() or query in value.lower() or query in fmt.uid):
            continue
        items.append(
            item(
                title=value,
                subtitle=fmt.label,
                arg=value,
                uid=fmt.uid,
            )
        )

    if not items:
        output(
            [
                item(
                    title=f'No format matches "{args}"',
                    subtitle="Try: YYYY, MM, DD, unix, ISO, -2d, +1w, 2026/7/9 ...",
                    valid=False,
                )
            ]
        )
        return

    output(items)
