### Why webhooks?

While polling and webhooks both accomplish the same task, webhooks are far more efficient. Zapier found that over 98.5% of polls are wasted. In contrast, webhooks only transfer data when there is new data to send, making them 100% efficient. That means that polling creates, on average, 66x more server load than webhooks. That’s a lot of wasted time, and if you’re paying per API call, a whole lot of wasted money.

When using polling, the frequency of your polls limits how up-to-date your event data is. For example, if your polling frequency is every 12 hours, the events returned by any poll could have happened any time in the past 12 hours. This means that any time an event occurs in the endpoint, your app will be out-of-date until the next poll.

With webhooks, this problem is eliminated. Since events are posted immediately to your monitored url, your apps will automatically update themselves with the new data almost instantly.

### Telegram instruction

