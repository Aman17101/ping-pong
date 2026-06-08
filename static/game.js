const canvas = document.getElementById('game');
const ctx = canvas.getContext('2d');
const socket = new WebSocket('ws://' + window.location.host + '/ws');

let paddleY = 150;

// Track mouse movement to move the paddle
canvas.addEventListener('mousemove', (e) => {
    const rect = canvas.getBoundingClientRect();
    paddleY = e.clientY - rect.top - 50; // Center the 100px paddle on the cursor
    socket.send(JSON.stringify({ y: paddleY }));
});

socket.onmessage = (event) => {
    const state = JSON.parse(event.data);
    draw(state);
};

const p1Display = document.getElementById('p1-score');
const p2Display = document.getElementById('p2-score');

function draw(state) {
    // 1. Clear the canvas
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    
    // 2. Update the Score HTML elements
    p1Display.innerText = state.Score1;
    p2Display.innerText = state.Score2;
    
    // 3. Check for a Winner
    if (state.Winner !== 0) {
        ctx.fillStyle = "white";
        ctx.font = "bold 40px 'Courier New'";
        ctx.textAlign = "center";
        ctx.fillText("PLAYER " + state.Winner + " WINS!", canvas.width / 2, canvas.height / 2);
        
        ctx.font = "20px 'Courier New'";
        ctx.fillText("Press F5 to Restart", canvas.width / 2, canvas.height / 2 + 50);
        
        // We still draw the paddles and ball, but they will be frozen 
        // because the Go backend stops updating them.
    }
    
    // 4. Draw Center Line
    ctx.setLineDash([5, 15]);
    ctx.beginPath();
    ctx.moveTo(400, 0);
    ctx.lineTo(400, 400);
    ctx.strokeStyle = "rgba(255, 255, 255, 0.5)";
    ctx.stroke();

    // 5. Draw Paddles
    ctx.fillStyle = 'white';
    ctx.setLineDash([]); // Ensure lines are solid
    // Player 1 (Left)
    ctx.fillRect(10, state.Paddle1Y, 10, 100);
    // Player 2 (Right)
    ctx.fillRect(780, state.Paddle2Y, 10, 100);
    
    // 6. Draw Ball
    ctx.beginPath();
    ctx.arc(state.BallX, state.BallY, 8, 0, Math.PI * 2);
    ctx.fill();
}