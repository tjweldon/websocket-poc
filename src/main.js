const flock = [];

let alignSlider, cohesionSlider, separationSlider;

const WIDTH = 1200;
const HEIGHT = 600;
const SLIDER_ROW = HEIGHT - 40;
const SLIDER_ROW_MARGIN = 40;
const SLIDER_SEP = 40;


const setup = (s) => () => {
    s.createCanvas(WIDTH, HEIGHT);
    s.background(0);

    var slider_cursor = { x: SLIDER_ROW_MARGIN, y: SLIDER_ROW };

    separationSlider = s.createSlider(0, 5, 0, 0.1);
    separationSlider.position(slider_cursor.x, slider_cursor.y);
    slider_cursor.x += separationSlider.width + SLIDER_SEP;

    cohesionSlider = s.createSlider(0, 5, 0, 0.1);
    cohesionSlider.position(slider_cursor.x, slider_cursor.y);
    slider_cursor.x += separationSlider.width + SLIDER_SEP;

    alignSlider = s.createSlider(0, 5, 0, 0.1);
    alignSlider.position(slider_cursor.x, slider_cursor.y);
    slider_cursor.x += separationSlider.width + SLIDER_SEP;

    for (let i = 0; i < 30; i++) {
        flock.push(new Boid(s));
    }
};

let i = 0;
const draw = (s) => () => {
    s.background(0);

    s.strokeWeight(3);
    s.textSize(14);
    s.text(`Separation: ${separationSlider.value()}`, separationSlider.x + 20, separationSlider.y - 10);
    s.text(`Cohesion: ${cohesionSlider.value()}`, cohesionSlider.x + 20, cohesionSlider.y - 10);
    s.text(`Align: ${alignSlider.value()}`, alignSlider.x + 20, alignSlider.y - 10);
    s.text(`BOID: ${flock[0]}`, 0, 50);

    for (let boid of flock) {
        boid.flock(s, flock);
        boid.update(s);
        boid.edges(s);
        boid.show(s);
    }
};

let sketch = (s) => {
    s.setup = setup(s);
    s.draw = draw(s);
};

const sketchInstance = new p5(sketch);