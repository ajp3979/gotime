# Go Time
Go Time is a Go app to help with time conversion and global timezones for the SOC. Go Time will also make a request to PagerDuty and return the on-call SIRT engineer given the provided schedule in the environment variables.

## Dependencies
Requires environment variables:
- `PAGER_DUTY_API_KEY` API key for PagerDuty. Read-only scope is recommended.
- `PAGER_DUTY_SCHEDULE` PagerDuty schedule ID.

## Example
```
$ gotime
Hyderabad: Tue Dec 17 09:35:02
UTC: Tue Dec 17 04:05:02
Dublin: Tue Dec 17 04:05:02
Eastern USA: Mon Dec 16 23:05:02
Central USA: Mon Dec 16 22:05:02
Pacific USA: Mon Dec 16 20:05:02
On-Call SIRT user@example.com
```