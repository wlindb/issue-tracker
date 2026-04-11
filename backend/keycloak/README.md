# Keycloak – Clerk-Inspired Login Theme

A custom Keycloak login theme that closely mirrors the clean, minimal aesthetic of [Clerk](https://clerk.com).

## Preview

- Centered white card with soft shadow on a light-gray background
- Inter font
- Violet (`#6c47ff`) primary button and focus ring
- Social-provider buttons (one per configured IdP)
- Inline field validation errors
- Password-visibility toggle
- "Sign up" footer link

## Theme structure

```
backend/keycloak/
├── realm-export.json
├── README.md
└── themes/
    └── clerk-inspired/
        └── login/
            ├── theme.properties          # parent theme, CSS reference
            ├── login.ftl                 # main login template
            └── resources/
                └── css/
                    └── login.css         # all styling
```

## Setup

### 1. Copy the theme into your Keycloak installation

#### Option A – Docker / Docker Compose (recommended)

Mount the theme folder as a volume into the Keycloak container:

```yaml
# docker-compose.yml (excerpt)
services:
  keycloak:
    image: quay.io/keycloak/keycloak:24.0
    volumes:
      - ./backend/keycloak/themes/clerk-inspired:/opt/keycloak/themes/clerk-inspired:ro
    command: start-dev
```

#### Option B – Manual copy

```bash
cp -r backend/keycloak/themes/clerk-inspired \
      /opt/keycloak/themes/clerk-inspired
```

#### Option C – Keycloak Operator (Kubernetes)

Use a Kubernetes `ConfigMap` or `PersistentVolume` to mount the theme directory into `/opt/keycloak/themes/`.

### 2. Activate the theme in the Keycloak Admin Console

1. Open **Keycloak Admin Console** → select your realm.
2. Go to **Realm Settings → Themes**.
3. Set **Login Theme** to `clerk-inspired`.
4. Click **Save**.

### 3. Verify

Navigate to your application's login URL (or click *Not logged in? Sign in* from the Admin Console preview). You should see the Clerk-inspired card UI.

## Customisation

| Variable | Default | Purpose |
|---|---|---|
| `--cl-primary` | `#6c47ff` | Button / focus / link colour |
| `--cl-bg` | `#f9fafb` | Page background |
| `--cl-card-bg` | `#ffffff` | Card background |
| `--cl-border` | `#e5e7eb` | Input / card border |
| `--cl-radius-lg` | `16px` | Card corner radius |
| `--cl-font-family` | Inter, system-ui | Body font |

All variables are defined at the top of `resources/css/login.css` under `:root`.

### Replacing the logo

The SVG logo in `login.ftl` is a placeholder. Replace the `<svg>` inside `.cl-logo-wrap` with:

- An `<img>` tag pointing to a file you place under `resources/img/`, **or**
- Your own inline SVG.

### Adding more pages

Keycloak has additional FreeMarker templates (`register.ftl`, `error.ftl`, `info.ftl`, etc.). Copy them from the base `keycloak` theme and apply the same card wrapper pattern used in `login.ftl`.

## Keycloak compatibility

Tested against **Keycloak 22 – 24** (Quarkus-based distribution).  
The theme inherits from the built-in `keycloak` parent theme, so all Keycloak-managed form actions, CSRF tokens, and social-login redirects are fully preserved.
