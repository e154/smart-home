The cron expression is made of six fields. Each field can have the following values.

| *             | *             | *             | *                         | *              | *                       |
|---------------|---------------|---------------|---------------------------|----------------|-------------------------|
| minute (0-59) | second (0-59) | hour (0 - 23) | day of the month (1 - 31) | month (1 - 12) | day of the week (0 - 6) |

Here are some examples for you.

| Cron expression | Schedule                              |
|-----------------|---------------------------------------|
| * * * * * *     | Every second                          |
| 0 * * * * *     | 	Every minute                         |
| 0 0 * * * *     | 	Every hour                           |
| 0 0 0 * * *     | 	Every day at 12:00 AM                |
| 0 0 0 * * FRI   | At 12:00 AM, only on Friday           |
| 0 0 0 1 * *     | 	At 12:00 AM, on day 1 of the month   |
