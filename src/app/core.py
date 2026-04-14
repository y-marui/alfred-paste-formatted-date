"""Application orchestrator.

Wires together the Router and Command handlers.
This is the single entry point called by scripts/entry.py.

To add a new format:
  Add a DateFormat entry to src/app/commands/date_cmd.py.
"""

from __future__ import annotations

from alfred.router import Router
from app.commands import config_cmd, date_cmd, help_cmd

router = Router(default="date")
router.register("date")(date_cmd.handle)
router.register("config")(config_cmd.handle)
router.register("help")(help_cmd.handle)


def run(query: str) -> None:
    """Main application entry point.

    Args:
        query: Raw query string from Alfred (e.g. "YYYY" to filter formats).
    """
    router.dispatch(query)
