class Boid {
    constructor(s) {
        this.position = s.createVector(
            s.width * Math.random(),
            s.height * Math.random()
        );
        this.velocity = p5.Vector.random2D();
        this.velocity.setMag((Math.random() + 0.5) * 5);
        this.acceleration = s.createVector();
        this.maxForce = 0.2;
        this.maxSpeed = 4;
    }

    edges(s) {
        if (this.position.x > s.width) {
            this.position.x = 0;
        } else if (this.position.x < 0) {
            this.position.x = s.width;
        }

        if (this.position.y > s.height) {
            this.position.y = 0;
        } else if (this.position.y < 0) {
            this.position.y = s.height;
        }
    }

    align(s, boids) {
        let perceptionRadius = 50;
        let steering = s.createVector();
        let total = 0;
        for (let neighbor of boids) {
            let d = s.dist(
                this.position.x,
                this.position.y,
                neighbor.position.x,
                neighbor.position.y
            );
            if (neighbor !== this && d < perceptionRadius) {
                steering.add(neighbor.velocity);
                total++;
            }
        }
        if (total > 0) {
            steering.div(total);
            steering.setMag(this.maxSpeed);
            steering.sub(this.velocity);
            steering.limit(this.maxForce);
        }

        return steering;
    }

    cohesion(s, boids) {
        let perceptionRadius = 50;
        let steering = s.createVector();
        let total = 0;
        for (let neighbor of boids) {
            let d = s.dist(
                this.position.x,
                this.position.y,
                neighbor.position.x,
                neighbor.position.y
            );
            if (neighbor !== this && d < perceptionRadius) {
                steering.add(neighbor.position);
                total++;
            }
        }
        if (total > 0) {
            steering.div(total);
            steering.sub(this.position);
            steering.setMag(this.maxSpeed);
            steering.sub(this.velocity);
            steering.limit(this.maxForce);
        }

        return steering;
    }

    separation(s, boids) {
        let perceptionRadius = 50;
        let steering = s.createVector();
        let total = 0;
        for (let neighbor of boids) {
            let d = s.dist(
                this.position.x,
                this.position.y,
                neighbor.position.x,
                neighbor.position.y
            );
            if (neighbor !== this && d < perceptionRadius) {
                let diff = p5.Vector.sub(this.position, neighbor.position);
                diff.div(d);
                steering.add(diff);
                total++;
            }
        }
        if (total > 0) {
            steering.div(total);
            steering.setMag(this.maxSpeed);
            steering.sub(this.velocity);
            steering.limit(this.maxForce);
        }

        return steering;
    }

    flock(s, boids) {
        this.acceleration.set(0, 0);
        let alignment = this.align(s, boids);
        let cohesion = this.cohesion(s, boids);
        let separation = this.separation(s, boids);

        separation.mult(separationSlider.value());
        cohesion.mult(cohesionSlider.value());
        alignment.mult(alignSlider.value());

        this.acceleration.add(alignment).add(cohesion).add(separation);
    }

    update(s) {
        if (s.dist(this.position.x, this.position.y, 600, 400) < 200) {
            if (this.velocity.x > this.velocity.y) {
                var vel = s.createVector(0, this.velocity.y)
            } else {
                var vel = s.createVector(this.velocity.y, 0)
            }
        } else {
            var vel = this.velocity
        }
        this.position.add(this.velocity);
        this.velocity.add(this.acceleration);
        this.velocity.limit(this.maxSpeed);
    }

    show(s) {
        s.strokeWeight(8);
        s.stroke(255);
        s.point(this.position.x, this.position.y);
    }

    toStr() {
        return JSON.stringify({"position": this.position, veldocity: this.velocity})
    }

}
