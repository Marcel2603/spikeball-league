# Architecture

## Overview

A lightweight, auth-less web application for tracking Spikeball leagues, teams, and match scores.
The system prioritizes extreme ease-of-use via a "StrawPoll/Doodle" creation model, eliminating user
registration while maintaining secure administrative control.

## Techstack

- Backend: Go (Golang)
- Frontend: HTMX (hypermedia-driven interactions), Alpine.js (client-side state & UI tweaks)
- Styling: Tailwind CSS
- Templates: a-h/templ (type-safe HTML generation)
- Database: SQLite (embedded, single-file database) with support for PostgreSQL
- Deployment: Docker

## Monitoring & Health

- **Metrics**: The application exports Prometheus-compatible metrics.
- **Health Checks**: Standard `/health/live` and `/health/ready` endpoints for container orchestration.

## AuthLess Mode

- To eliminate login friction, the application uses a token-based URL routing architecture rather than traditional sessions.
- Public Link (/l/<public_id>): Shared with the friend group. Used to view standings, match history, and submit new
game scores.
- Admin Link (/admin/<admin_id>): Kept strictly by the creator. Used to manage the league (rename, delete matches,
edit players).
- Client-Side Persistence: When an admin creates a league, Alpine.js intercepts the response and stores the admin_id
and League Name
in the browser's localStorage. The index page (/) reads this to provide a "Your Recent Leagues" shortcut for returning creators.

## Static Assets

The documentation and certain frontend assets are pre-generated using internal tools and bundled or served alongside the
Go binary to ensure a single-dependency deployment model.
