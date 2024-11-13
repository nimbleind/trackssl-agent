# TrackSSL Agent

This agent is installed on your local network to facilitate monitoring of
certificates that are not available on the public internet. To utilize, you'll
need a TrackSSL account on a plan with Private Certificate Monitoring. You'll
also need a TrackSSL API key and an Agent Token and you'll need to assign
domains to agents in the TrackSSL dashboard.

For help docs, please see [TrackSSL Documentation](https://trackssl.com/help/).

## Installation

You can download a precompiled binary from the
[releases page](https://github.com/trackssl/trackssl-agent-go/releases).

## Execution

The agent requires runs as a daemon, sleeping by default for 4 hours between checks.

### Windows

Simply run it with your agent token and API key:

```
$ .\trackssl-agent-windows.exe -agent-token your_agent_token -auth-token your_api_key
```

### Linux

You will need to make the binary executable before running it:

```
$ chmod +x trackssl-agent-linux
```

Then run it with your agent token and API key:

```
$ ./trackssl-agent-linux -agent-token your_agent_token -auth-token your_api_key
```

### Mac

You will need to make the binary executable before running it:

```
$ chmod +x trackssl-agent-mac
```

And then approve the unsigned binary via the Security & Privacy settings in System Preferences.

Then run it with your agent token and API key:

```
$ ./trackssl-agent-mac -agent-token your_agent_token -auth-token your_api_key
```

## Execution at Start Up

### Windows

To add the TrackSSL Agent as a startup task on Windows using Task Scheduler, follow these steps:

1. **Open Task Scheduler**:
   - Press `Win + S`, type "Task Scheduler," and open it.

2. **Create a New Task**:
   - In Task Scheduler, select **Action > Create Task**.

3. **Configure General Settings**:
   - In the **General** tab, name the task, e.g., "TrackSSL Agent".
   - Set the task to run **only when the user is logged on** or **whether the user is logged on or not**, depending on your preference.

4. **Set the Trigger**:
   - Go to the **Triggers** tab and click **New**.
   - Select **At log on** or **At startup** as the trigger type.
   - Click **OK** to save the trigger.

5. **Configure the Action**:
   - Go to the **Actions** tab and click **New**.
   - Select **Start a program** as the action.
   - In the **Program/script** field, browse to the location of `trackssl-agent-windows.exe`.
   - In the **Add arguments** field, enter your TrackSSL agent token and API key, like so:
     ```
     -agent-token your_agent_token -auth-token your_api_key
     ```
   - Click **OK** to save the action.

6. **Finish and Save**:
   - Review the settings and click **OK** to save the task.

Your TrackSSL Agent will now automatically start each time your computer starts up or when you log on. For more information or troubleshooting, visit the [TrackSSL Documentation](https://trackssl.com/help/).

## Building

To compile the agent for all platforms, run the following command:

```
$ make
```
