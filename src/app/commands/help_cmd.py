"""help command - display available commands.

Usage in Alfred:  date help
"""

from __future__ import annotations

from alfred.response import item, output

_COMMANDS = [
    ("date", "List all date formats (default command)", "date "),
    ("date config", "View or reset configuration", "date config"),
    ("date help", "Show this help", "date help"),
]


def handle(args: str) -> None:  # noqa: ARG001
    """Display all available commands."""
    output(
        [
            item(
                title=cmd,
                subtitle=desc,
                arg="",
                uid=f"help-{cmd.split()[0]}",
                valid=False,
                autocomplete=autocomplete,
            )
            for cmd, desc, autocomplete in _COMMANDS
        ]
    )
