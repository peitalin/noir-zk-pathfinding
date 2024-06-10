use dep::std;
use dep::std::ec::tecurve::affine::{Curve, Point};

// same function as std::merkle::compute_merkle_root() when using pedersen_hash
fn compute_root_from_leaf(leaf : Field, index : Field, hash_path: [Field; 2]) -> Field {
    let n = hash_path.len();
    let index_bits = index.to_le_bits(n as u32);
    println(f"index bits: {index_bits}");
    let mut current = leaf;
    for i in 0..n {
        let path_bit = index_bits[i] as bool;
        println(f"path bit: {path_bit}");
        let (hash_left, hash_right) = if path_bit {
            (hash_path[i], current)
        } else {
            (current, hash_path[i])
        };
        // Returns the root of the tree from the provided leaf and its hashpath, using mimc hash
        // current = std::hash::mimc::mimc_bn254([hash_left, hash_right]);
        current = std::hash::pedersen_hash([hash_left, hash_right]);
        println(f"current hash: {current}");
    };
    current
}