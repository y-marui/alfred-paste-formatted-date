"""date command - list date formats for selection and paste.

Usage in Alfred:  date
                  date <filter>   — filter by format name or example
"""

from __future__ import annotations

import datetime
from typing import NamedTuple

from alfred.logger import get_logger
from alfred.response import item, output

log = get_logger(__name__)


class DateFormat(NamedTuple):
    uid: str
    label: str
    strftime: str | None
    custom: str | None = None  # used when strftime is not sufficient


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
    """Show date formats; optionally filter by format label or example value."""
    log.debug("date command: args=%r", args)

    now = datetime.datetime.now()
    query = args.strip().lower()

    items = []
    for fmt in _FORMATS:
        value = _formatted(fmt, now)
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
                    subtitle="Try: YYYY, MM, DD, unix, ISO ...",
                    valid=False,
                )
            ]
        )
        return

    output(items)
