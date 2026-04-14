"""Tests for command handlers."""

from __future__ import annotations

import json

from app.commands import config_cmd, date_cmd, help_cmd


class TestDateCommand:
    def test_empty_query_returns_all_formats(self, capsys):
        date_cmd.handle("")
        data = json.loads(capsys.readouterr().out)
        assert len(data["items"]) == len(date_cmd._FORMATS)

    def test_all_items_are_valid(self, capsys):
        date_cmd.handle("")
        data = json.loads(capsys.readouterr().out)
        assert all(it["valid"] for it in data["items"])

    def test_filter_by_format_label(self, capsys):
        date_cmd.handle("ISO")
        data = json.loads(capsys.readouterr().out)
        assert len(data["items"]) > 0
        # Each result must match "iso" in its label, value, or uid
        for it in data["items"]:
            matches = (
                "iso" in it["subtitle"].lower()
                or "iso" in it["title"].lower()
                or "iso" in it["uid"].lower()
            )
            assert matches

    def test_filter_by_format_uid_yyyymmdd(self, capsys):
        date_cmd.handle("YYYYMMDD")
        data = json.loads(capsys.readouterr().out)
        assert len(data["items"]) >= 1
        assert data["items"][0]["subtitle"] == "YYYYMMDD"

    def test_no_match_returns_error_item(self, capsys):
        date_cmd.handle("xyzzy-nonexistent")
        data = json.loads(capsys.readouterr().out)
        assert len(data["items"]) == 1
        assert data["items"][0]["valid"] is False

    def test_arg_equals_title(self, capsys):
        """The arg (value to paste) must equal the displayed date string."""
        date_cmd.handle("")
        data = json.loads(capsys.readouterr().out)
        for it in data["items"]:
            assert it["arg"] == it["title"]

    def test_unix_timestamp_is_numeric(self, capsys):
        date_cmd.handle("unix")
        data = json.loads(capsys.readouterr().out)
        assert len(data["items"]) == 1
        assert data["items"][0]["arg"].isdigit()


class TestConfigCommand:
    def test_empty_config_shows_no_settings(self, capsys):
        config_cmd.handle("")
        data = json.loads(capsys.readouterr().out)
        titles = [it["title"] for it in data["items"]]
        assert any("No settings" in t for t in titles)

    def test_reset_clears_config(self, capsys):
        config_cmd._config.set("key", "value")
        config_cmd.handle("reset")
        data = json.loads(capsys.readouterr().out)
        assert "reset" in data["items"][0]["title"].lower()
        assert config_cmd._config.all() == {}

    def test_shows_existing_settings(self, capsys):
        config_cmd._config.set("some_key", "some_value")
        config_cmd.handle("")
        data = json.loads(capsys.readouterr().out)
        titles = [it["title"] for it in data["items"]]
        assert any("some_key" in t for t in titles)

    def test_unknown_subcommand_shows_current_config(self, capsys):
        config_cmd.handle("unknown-subcommand")
        data = json.loads(capsys.readouterr().out)
        assert len(data["items"]) > 0


class TestHelpCommand:
    def test_shows_all_commands(self, capsys):
        help_cmd.handle("")
        data = json.loads(capsys.readouterr().out)
        assert len(data["items"]) == len(help_cmd._COMMANDS)

    def test_all_items_invalid(self, capsys):
        help_cmd.handle("")
        data = json.loads(capsys.readouterr().out)
        assert all(not it["valid"] for it in data["items"])
