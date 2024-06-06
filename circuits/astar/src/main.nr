mod foo;
use foo::bar;
use dep::std::slice;


fn main(start: [i32; 2], end: [i32; 2], max_steps: pub i32) -> pub i32 {

    let x1 = start[0];
    let y1 = start[1];

    let x2 = end[0];
    let y2 = end[1];

    let astar_path_results = get_astar_path([x1, y1], [x2, y2]);
    let path_pruned = pruneAstarPath(astar_path_results);
    let path_len = path_pruned.len() as i32;
    /// MUST REMOVE println for noir-js proof generation to work ////
    // println(f"A* path: {path_pruned}");
    // println(f"assert A* path length: {path_len} <= max_steps: {max_steps}");
    assert(path_len <= max_steps);

    path_len
}

#[oracle(GetAstarPath)]
unconstrained fn astar_path_oracle(start: [i32; 2], end: [i32; 2]) -> (
    // max number of path points returned from the orcacle is static
    [i32; 2],
    [i32; 2],
    [i32; 2],
    [i32; 2],
    [i32; 2],
    [i32; 2],
    [i32; 2],
    [i32; 2],
    [i32; 2],
    [i32; 2],
) {}

unconstrained fn get_astar_path(
    start: [i32; 2],
    end: [i32; 2]
) -> [[i32;2]; 10] {
    let (r0, r1, r2, r3, r4, r5, r6, r7, r8, r9) = astar_path_oracle(start, end);
    let mut path = [r0, r1, r2, r3, r4, r5, r6, r7, r8, r9];
    path
}

fn pruneAstarPath(pathArray: [[i32; 2]; 10]) -> [[i32; 2]] {

    let mut path = pathArray.as_slice();
    let NOIR_ZERO_VALUE = 65_536;
    // Use this value as "zero" when pruning values
    // see NOIR_ZERO_VALUE_HEX in the oracle/cmd/main.go
    // hex value: 0x10000 => 65536 when decoded as i32.
    for i in 0..9 {
        let lastIndex = 9 - i;
        if path[lastIndex][0] == NOIR_ZERO_VALUE {
            if path[lastIndex][1] == NOIR_ZERO_VALUE {
                let res = path.remove(lastIndex);
                path = res.0;
            }
        }
    }
    path
}


#[test]
fn test_main() {
    let max_steps = 10;
    let steps = main([1, 1], [4, 5], max_steps);
    println(f"A* number of steps: {steps}");
}

#[test]
fn test_astar_oracle1() {
    let mut path = get_astar_path([0, 0], [7, 2]);
    // Noir oracles can only receive data in Single or Array format
    // Arrays must specify specific lengths,
    // so we return zero'd out path coordinates as [NOIR_ZERO_VALUE, NOIR_ZERO_VALUE] and prune them
    let pathPruned = pruneAstarPath(path);
    println(f"pruned path: {pathPruned}");
}

#[test]
fn test_astar_oracle2() {
    let mut path = get_astar_path([0, 0], [9, 8]);
    let pathPruned = pruneAstarPath(path);
    println(f"pruned path: {pathPruned}");
}
