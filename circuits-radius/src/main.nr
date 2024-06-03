
fn main(x1: i32, x2: i32, y1: i32, y2: i32, d: pub i32) {

    // distance = sqrt( (x2 - x1)^2 + (y2 - y1)^2 )
    let dx = x2 - x1;
    let dx2 = dx * dx;

    let dy = y2 - y1;
    let dy2 = dy * dy;

    let dx2_dy2 = dx2 + dy2;

    // instead of sqrt, we square distance
    // let maxdist = d * d;
    let maxdist = d;

    let sqrtd = get_sqrt(dx2_dy2);
    // println(f"ceil(sqrt) of {dx2_dy2} = {sqrtd}"); // remove println for noir-js
    // let sqrtd = dx2_dy2;

    assert(sqrtd <= maxdist);

}


#[oracle(GetSqrt)]
unconstrained fn sqrt(d: i32) -> i32 { }

unconstrained fn get_sqrt(d: i32) -> i32 {
    sqrt(d)
}


#[test]
fn test_main() {
    main(1, 4, 1, 5, 49);

}

#[test]
fn test_main_oracle() {
    let maxdist = get_sqrt(225);
    main(1, 4, 1, 6, maxdist);

}


