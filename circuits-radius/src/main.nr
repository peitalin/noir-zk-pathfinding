
fn main(x1: u32, x2: u32, y1: u32, y2: u32, d: pub u32) {

    // distance = sqrt( (x2 - x1)^2 + (y2 - y1)^2 )
    let dx = x2 - x1;
    let dx2 = dx * dx;

    let dy = y2 - y1;
    let dy2 = dy * dy;

    let dx2_dy2 = dx2 + dy2;

    // instead of sqrt, we square distance
    let maxdist = d * d;

    assert(dx2_dy2 <= maxdist);
}

#[test]
fn test_main() {
    main(1, 2, 3, 4, 50);

}

