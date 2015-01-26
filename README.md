SendGrid Webhook Handler for CRM Bliss
======================================

Testing with curl:

```
curl -i -X POST http://localhost:3000/ -H "Content-Type: application/json" \
  -d '{"email_event_id":"f69e09b9-0617-467a-87e6-90d9851ce538", "event":"something!"}'
```
