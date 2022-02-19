const flock = [];
let frameBuffer = []

const WIDTH = 1200;
const HEIGHT = 600;
const SLIDER_ROW = HEIGHT - 40;
const SLIDER_ROW_MARGIN = 40;
const SLIDER_SEP = 40;


const setup = (s) => () => {
    s.createCanvas(WIDTH, HEIGHT);
    s.background(0);
};


let i = 0;
const draw = (s) => () => {
    s.background(0);
    show(s);
};

let sketch = (s) => {
    s.setup = setup(s);
    s.draw = draw(s);
};

const show = (s) => {
    if (frameBuffer.length < 1) {
        return
    }
    let position;
    s.background(0)
    let frame = frameBuffer.shift();
    if (frame.includes("{")) {
        position = JSON.parse(frame);
    } else {
        return;
    }


    s.strokeWeight(8);
    s.stroke(255);
    s.point(position.x, position.y);
}


const sketchInstance = new p5(sketch);