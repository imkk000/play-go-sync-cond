# Play Sync.Cond

## How it work?

- Simulate request from client with random context timeout
- After some delay, sender send new message into a queue then broadcast message
- Waiter wake up then process message but before do it. It checks message still empty in the queue
- If not empty, break the loop then process
- If empty, check context is cancel? If not cancel, wait again. If cancel, return nil
- On exit, close done channel to signal waker goroutine
- Waker exits silently if message processed (no broadcast needed)
- Waker broadcasts only if context cancelled, to unblock wait
