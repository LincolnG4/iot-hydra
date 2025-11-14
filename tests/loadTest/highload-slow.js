import ws from 'k6/ws';
import { check } from 'k6';

// Reusable function for the core WebSocket logic
function websocketSession(sessionDuration) {
  const url = 'ws://localhost:8080/v1/ws';
  const payload = JSON.stringify({
    topic: 'my.iot',
    payload: 'SGVsbG8gV29ybGQ=', // "Hello World" in Base64
    target_brokers: ['ligmaNats'],
  });

  const res = ws.connect(url, {}, function (socket) {
    socket.on('open', () => {
      // Send a message as soon as the connection is established
      socket.send(payload);
    });

    socket.on('error', (e) => {
      console.log('An error occurred: ', e.error());
    });

    // Keep the connection open for the specified duration
    socket.setTimeout(() => {
      socket.close();
    }, sessionDuration);
  });

  check(res, { 'status is 101': (r) => r && r.status === 101 });
}

// ========================================================================
//                            SCENARIO DEFINITIONS
// ========================================================================
export const options = {
  // Define the different device behaviors as separate scenarios
  scenarios: {
    // Scenario 1: Stable Devices (Always-on)
    // Simulates a baseline of devices that are always connected.
    stable_connections: {
      executor: 'constant-vus',
      exec: 'stableDeviceFlow', // The function this scenario will run
      vus: 50, // Start and maintain 50 concurrent connections
      duration: '5m', // Run this scenario for 5 minutes
    },

    // Scenario 2: Intermittent Devices (Connect, Send, Disconnect)
    // Simulates devices that connect frequently but don't stay connected.
    intermittent_connections: {
      executor: 'ramping-arrival-rate',
      exec: 'intermittentDeviceFlow', // The function this scenario will run
      startRate: 10, // Start with 10 new connections per second
      timeUnit: '1s', // The period for the rate (e.g., 10 per 1s)
      preAllocatedVUs: 50, // VUs to pre-allocate for this scenario
      maxVUs: 100, // If the rate requires it, scale up to this many VUs
      stages: [
        { target: 20, duration: '1m' }, // Ramp up to 20 new connections/sec over 1 minute
        { target: 20, duration: '3m' }, // Maintain this rate for 3 minutes
        { target: 0, duration: '1m' },  // Ramp down
      ],
    },

    // Scenario 3: A Sudden Reconnection Spike
    // Simulates a power-recovery event where many devices reconnect at once.
    reconnection_spike: {
      executor: 'ramping-vus',
      exec: 'intermittentDeviceFlow', // These devices also behave like intermittent ones
      startTime: '2m', // Start this spike 2 minutes into the main test
      stages: [
        { target: 150, duration: '10s' }, // Quickly ramp up to 150 VUs
        { target: 150, duration: '30s' }, // Hold the spike for 30 seconds
        { target: 0, duration: '10s' },   // Ramp down
      ],
    },
  },
  // Set thresholds based on your SLOs. The test will fail if these are breached.
  thresholds: {
    'checks': ['rate>0.999'], // 99.9% of connections must succeed
    'ws_connecting': ['p(95)<200'], // 95% of connections must be faster than 200ms
  },
};

// ========================================================================
//                      EXECUTION FUNCTIONS FOR SCENARIOS
// ========================================================================

// Logic for a stable device: connect and stay connected for a long time
export function stableDeviceFlow() {
  // Stay connected for a long time, e.g., ~1 minute for the test's purpose
  const long_duration = 60000;
  websocketSession(long_duration);
}

// Logic for an intermittent device: connect briefly and disconnect
export function intermittentDeviceFlow() {
  // Stay connected for a short, random time
  const short_duration = Math.random() * 2000 + 1000; // 1s to 3s
  websocketSession(short_duration);
}
