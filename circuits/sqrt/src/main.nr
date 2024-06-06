mod foo;
use foo::bar;
use dep::std::slice;


fn main(start: [i32; 2], end: [i32; 2], max_dist: pub i32) -> pub i32 {

    let x1 = start[0];
    let y1 = start[1];

    let x2 = end[0];
    let y2 = end[1];

    // Prove move distance doesn't exceed max_dist
    // distance = sqrt( (x2 - x1)^2 + (y2 - y1)^2 )
    let dx = x2 - x1;
    let dx2 = dx * dx;
    let dy = y2 - y1;
    let dy2 = dy * dy;

    let dx2_dy2 = dx2 + dy2;
    let distance = get_sqrt( dx2_dy2 ) as i32;
    let _sqrt_sq = distance * distance;
    assert(_sqrt_sq == dx2_dy2); // constrain the oracle response

    // println(f"ceil(sqrt) of {dx2_dy2} = {sqrtd}"); // remove println for noir-js
    // check move distance doesn't exceed max_dist
    assert(distance <= max_dist);

    distance
}


#[oracle(GetSqrt)]
unconstrained fn sqrt(d: i32) -> i32 { }

unconstrained fn get_sqrt(d: i32) -> i32 {
    sqrt(d)
}

#[test]
fn test_main() {
    let max_distance = 49;
    let distance = main([1, 1], [4, 5], max_distance);
    println(f"Distance: {distance}");
}

#[test]
fn test_sqrt_oracle() {
    let sqrtnum = 225;
    let maxdist = get_sqrt(sqrtnum);
    println(f"sqrt({sqrtnum}) = {maxdist}");
    assert(maxdist == 15);
}
