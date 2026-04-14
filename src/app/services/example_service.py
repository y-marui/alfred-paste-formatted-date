"""ExampleService - replace this with your own service.

A service encapsulates business logic and coordinates between
commands and external clients (APIs, local data, etc.).
"""

from __future__ import annotations

import os
from typing import Any, cast

from alfred.cache import Cache
from alfred.logger import get_logger
from app.clients.api_client import ApiClient

log = get_logger(__name__)


class ExampleService:
    """Fetches and caches search results.

    Replace this with your own service class that wraps whatever
    data source your workflow targets.
    """

    def __init__(self, ttl: int | None = None) -> None:
        resolved_ttl = ttl if ttl is not None else int(os.environ.get("cache_ttl", "300"))
        self._cache = Cache(ttl=resolved_ttl, namespace="example_service")
        self._client = ApiClient()

    def search(self, query: str) -> list[dict[str, Any]]:
        """Return search results for query, using cache when available.

        Args:
            query: Search string.

        Returns:
            List of result dicts with keys: id, title, subtitle, url.
        """
        log.debug("search: query=%r", query)

        cached = self._cache.get(query)
        if cached is not None:
            log.debug("cache hit: query=%r", query)
            return cast(list[dict[str, Any]], cached)

        results = self._client.search(query)
        self._cache.set(query, results)
        return results
