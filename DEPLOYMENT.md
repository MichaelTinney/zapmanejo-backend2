# DigitalOcean App Platform Deployment Guide

## Pre-Deployment Checklist

Before deploying to DigitalOcean App Platform, ensure you've completed these critical steps:

### 1. Update app.yaml with Your GitHub Repository

**File:** `app.yaml` (line 6)

Replace the placeholder with your actual GitHub username:

```yaml
repo_clone_url: https://github.com/YOUR_ACTUAL_USERNAME/zapmanejo-cleanbackend.git
```

### 2. Generate Strong Secrets

You need to generate secure random values for:

- **JWT_SECRET** - For signing authentication tokens
- **WHATSAPP_VERIFY_TOKEN** - For WhatsApp webhook verification

**Generate secrets using:**

```bash
# Generate JWT_SECRET (32+ characters)
openssl rand -base64 32

# Generate WHATSAPP_VERIFY_TOKEN
openssl rand -base64 24
```

Save these values - you'll need them in the next step.

---

## Deployment Steps

### Step 1: Set Up Environment Variables in DigitalOcean Console

**CRITICAL:** Do NOT add secrets to `app.yaml`. Set them in the DigitalOcean Console instead.

1. Go to: **DigitalOcean Dashboard** > **Apps** > **Your App** > **Settings**
2. Navigate to: **Components** > **backend** > **Environment Variables**
3. Add the following as **encrypted** environment variables:

| Variable Name | Value Source | Required |
|--------------|--------------|----------|
| `DATABASE_URL` | From your DO PostgreSQL database "Connection Details" | ✅ YES |
| `JWT_SECRET` | Generated value from Step 2 above | ✅ YES |
| `WHATSAPP_VERIFY_TOKEN` | Generated value from Step 2 above | ✅ YES |
| `WHATSAPP_NUMBER` | Your WhatsApp Business phone number (no + or spaces) | ✅ YES |
| `ALLOWED_ORIGINS` | Your frontend URL (e.g., `https://zapmanejo.com`) | ⚠️ Recommended |

**Notes:**
- `PORT` and `LIFETIME_TOTAL` are already set in `app.yaml` - no need to add them
- Mark all sensitive variables as **encrypted** in the DO console
- Use the "Connection String" format for `DATABASE_URL`, not individual fields

### Step 2: Get Your Database Connection String

1. Go to: **DigitalOcean Dashboard** > **Databases** > **Your PostgreSQL Database**
2. Click: **Connection Details**
3. Select: **Connection String** (not Connection Parameters)
4. Copy the full connection string (should look like):
   ```
   postgres://username:password@host:25060/defaultdb?sslmode=require
   ```
5. Use this as the value for `DATABASE_URL` in Step 1

### Step 3: Deploy to DigitalOcean

#### Option A: Deploy via DigitalOcean Console

1. Go to: **DigitalOcean Dashboard** > **Apps**
2. Click: **Create App**
3. Select: **GitHub** as source
4. Choose your repository: `zapmanejo-cleanbackend`
5. Select branch: `main`
6. DigitalOcean will auto-detect the `app.yaml` configuration
7. Review settings and click: **Next**
8. Add environment variables (from Step 1)
9. Click: **Create Resources**

#### Option B: Deploy via DigitalOcean CLI (doctl)

```bash
# Install doctl if you haven't already
# https://docs.digitalocean.com/reference/doctl/how-to/install/

# Authenticate
doctl auth init

# Create app from app.yaml
doctl apps create --spec app.yaml

# Get app ID
doctl apps list

# Update environment variables (via console is easier for secrets)
```

### Step 4: Verify Deployment

1. Wait for deployment to complete (typically 3-5 minutes)
2. Get your app URL from DigitalOcean Console
3. Test the health endpoint:
   ```bash
   curl https://your-app-url.ondigitalocean.app/health
   ```
4. Expected response:
   ```json
   {
     "status": "healthy",
     "database": "connected",
     "app": "ZapManejo v1.0"
   }
   ```

### Step 5: Update CORS Settings

Once you have your frontend deployed:

1. Go back to: **App Settings** > **Environment Variables**
2. Update `ALLOWED_ORIGINS` with your actual frontend URL:
   ```
   https://your-frontend-domain.com
   ```
3. For multiple origins, use comma-separated values:
   ```
   https://zapmanejo.com,https://www.zapmanejo.com
   ```

---

## Post-Deployment Verification

### Test All Endpoints

```bash
# Set your app URL
APP_URL="https://your-app-url.ondigitalocean.app"

# 1. Test root endpoint
curl $APP_URL/

# 2. Test health endpoint
curl $APP_URL/health

# 3. Test user registration
curl -X POST $APP_URL/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test123!","name":"Test User"}'

# 4. Test user login
curl -X POST $APP_URL/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test123!"}'
```

### Check Logs

```bash
# Via doctl CLI
doctl apps logs <app-id> --follow

# Or via DigitalOcean Console:
# Dashboard > Apps > Your App > Runtime Logs
```

Look for these success messages:
```
✓ All required environment variables are set
✓ Connected to DigitalOcean PostgreSQL with connection pooling
Starting database migration...
✓ Database schema migrated successfully
✓ Indexes created successfully
✓ Migration completed successfully
ZapManejo backend running on :8080
```

---

## Troubleshooting

### Issue: "DATABASE_URL not set" Error

**Cause:** Environment variable not configured in DO Console

**Fix:**
1. Go to App Settings > Environment Variables
2. Add `DATABASE_URL` with your PostgreSQL connection string
3. Redeploy the app

### Issue: "Failed to connect to database" Error

**Cause:** Incorrect database connection string or database not accessible

**Fix:**
1. Verify the connection string in DO Console
2. Ensure your PostgreSQL database is running
3. Check that `sslmode=require` is included in the connection string
4. Verify the database allows connections from your app

### Issue: Health Check Failing (503 Service Unavailable)

**Cause:** Database connectivity issues

**Fix:**
1. Check app logs for specific error messages
2. Verify database is running and accessible
3. Test database connection string manually
4. Check connection pool settings aren't too restrictive

### Issue: "too many connections" Error

**Cause:** Connection pool limits exceeded

**Fix:**
This should not happen with the new connection pooling settings, but if it does:
1. Check your DigitalOcean PostgreSQL plan connection limit
2. Reduce `SetMaxOpenConns` in `internal/database/db.go:35`
3. Current setting: 20 connections (safe for basic plan)

### Issue: CORS Errors from Frontend

**Cause:** Frontend domain not in `ALLOWED_ORIGINS`

**Fix:**
1. Add your frontend URL to `ALLOWED_ORIGINS` environment variable
2. Format: `https://your-frontend.com` (no trailing slash)
3. For multiple: `https://app.com,https://www.app.com`

### Issue: JWT Token Invalid/Expired

**Cause:** `JWT_SECRET` not set or changed

**Fix:**
1. Ensure `JWT_SECRET` is set in DO Console
2. Never change JWT_SECRET in production (invalidates all tokens)
3. Generate a new secret only if necessary, knowing all users will need to re-login

---

## Scaling Considerations

### Current Configuration

- **Instance Count:** 1 (single instance)
- **Instance Size:** basic-xs (512MB RAM, 1 vCPU)
- **Database Connections:** Max 20 per instance

### When to Scale Up

**Increase Instance Size** when:
- CPU usage consistently above 80%
- Memory usage above 90%
- Response times degrading

**Increase Instance Count** when:
- Need high availability (2+ instances recommended)
- Traffic exceeds single instance capacity
- Want zero-downtime deployments

**To scale instances:**

1. Edit `app.yaml`:
   ```yaml
   instance_count: 2  # or more
   instance_size_slug: basic-s  # or larger
   ```

2. Adjust database connections per instance:
   ```go
   // In internal/database/db.go:35
   sqlDB.SetMaxOpenConns(10)  // Reduce per instance when scaling horizontally
   ```

3. **Important:** With multiple instances and connection limit of 22:
   - 2 instances: Max 10 connections each = 20 total (safe)
   - 3 instances: Max 6 connections each = 18 total (safe)

---

## Security Checklist

- ✅ All secrets removed from `app.yaml`
- ✅ Environment variables set as "encrypted" in DO Console
- ✅ JWT_SECRET is strong random value (32+ characters)
- ✅ WHATSAPP_VERIFY_TOKEN is strong random value
- ✅ CORS restricted to actual frontend domain
- ✅ Database uses SSL (`sslmode=require`)
- ✅ Health endpoint configured for monitoring
- ✅ Connection pooling prevents resource exhaustion
- ✅ `.env` file in `.gitignore` (check this!)

---

## Monitoring

### Built-in Health Checks

DigitalOcean will automatically monitor:
- HTTP endpoint: `/health`
- Check interval: Every 10 seconds
- Timeout: 3 seconds
- Failure threshold: 3 consecutive failures

If health checks fail, DO will:
1. Mark the instance as unhealthy
2. Stop routing traffic to it
3. Attempt to restart the container

### Custom Monitoring

You can add external monitoring:
- **UptimeRobot** - Free tier available
- **Pingdom** - HTTP monitoring
- **New Relic** - APM (application performance monitoring)

Monitor these endpoints:
- `GET /health` - Database connectivity
- `GET /` - Basic API availability
- `POST /api/auth/login` - Critical user flow

---

## Backup and Recovery

### Database Backups

DigitalOcean automatically backs up your PostgreSQL database:
- Daily backups retained for 7 days (basic plan)
- Point-in-time recovery available (professional plans)

**Manual backup:**
```bash
# Get database credentials from DO Console
pg_dump -h HOST -U USERNAME -d DATABASE > backup.sql
```

### Application Rollback

If deployment fails:

1. Via Console: Apps > Deployments > Click previous deployment > "Redeploy"
2. Via CLI:
   ```bash
   doctl apps list-deployments <app-id>
   doctl apps create-deployment <app-id> --deployment-id <previous-deployment-id>
   ```

---

## Support Resources

- **DigitalOcean Docs:** https://docs.digitalocean.com/products/app-platform/
- **Community:** https://www.digitalocean.com/community/
- **Status Page:** https://status.digitalocean.com/

---

## Quick Reference Commands

```bash
# View app logs
doctl apps logs <app-id> --follow

# List deployments
doctl apps list-deployments <app-id>

# Get app info
doctl apps get <app-id>

# Update from app.yaml
doctl apps update <app-id> --spec app.yaml

# Restart app
doctl apps create-deployment <app-id>
```

---

## Next Steps After Deployment

1. ✅ Test all API endpoints
2. ✅ Configure WhatsApp webhook with your app URL
3. ✅ Set up PayPal credentials (if not done)
4. ✅ Update frontend to use production API URL
5. ✅ Set up monitoring/alerting
6. ✅ Configure custom domain (optional)
7. ✅ Enable HTTPS redirect (DO does this automatically)
8. ✅ Set up error tracking (Sentry, Rollbar, etc.)

---

**Deployment Date:** _____________________
**App URL:** _____________________
**Database:** _____________________
**Version:** v1.0
