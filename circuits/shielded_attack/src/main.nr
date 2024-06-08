use dep::std;
use dep::std::ec::tecurve::affine::{Curve, Point};


fn main(
  recipient: pub Field,
  // Private key of note
  // all notes have the same denomination
  priv_key: Field,
  // Merkle membership proof
  note_root: pub Field,
  index: Field,
  note_hash_path: [Field; 3],
  // Random secret to keep note_commitment private
  secret: Field,
  // Hash to be checked against the nullifier computed in the circuit
  nullifierHash: pub Field,
) -> pub Field {
    // Compute public key from private key to show ownership
    let pubkey = std::embedded_curve_ops::fixed_base_scalar_mul(priv_key, 1);
    println(f"pubkey1: {pubkey}");

    // let pubkey2 = std::embedded_curve_ops::fixed_base_scalar_mul(1, priv_key);
    // println(f"pubkey2: {pubkey2}");

    let pubkey_x = pubkey[0];
    let pubkey_y = pubkey[1];

    // Compute input note commitment
    let note_commitment = std::hash::pedersen_hash([pubkey_x, pubkey_y, secret]);

    println("");
    // Compute input note nullifier
    let nullifier = std::hash::pedersen_hash([note_commitment, index, priv_key]);
    println(f"nullifier:     {nullifier}");
    println(f"nullifierHash: {nullifierHash}");
    let nullifierEq = nullifierHash == nullifier;
    println(f"nullifierHash == nullifier: {nullifierEq}");
    // assert(nullifierHash == nullifier);

    println("");
    // Check that the input note commitment is in the root
    let new_root = compute_root_from_leaf(note_commitment, index, note_hash_path);
    println(f"new_root:  {new_root}");
    println(f"note_root: {note_root}");
    // assert(new_root == note_root);

    // Cannot have unused variables, return the recipient as public output of the circuit
    recipient
}

// Returns the root of the tree from the provided leaf and its hashpath, using mimc hash
fn compute_root_from_leaf(leaf : Field, index : Field, hash_path: [Field; 3]) -> Field {
    let n = hash_path.len();
    let index_bits = index.to_le_bits(n as u32);
    let mut current = leaf;
    for i in 0..n {
        let path_bit = index_bits[i] as bool;
        let (hash_left, hash_right) = if path_bit {
            (hash_path[i], current)
        } else {
            (current, hash_path[i])
        };

        current = std::hash::mimc::mimc_bn254([hash_left, hash_right]);
    };
    current
}


#[test]
fn test_main() {
    // # Random test key
    let recipient = 0x00000000000000000000000070997970c51812dc3a010c7d01b50e0d17dc79c8;
    let priv_key = 0x000000000000000000000000000000000000000000000000000000616c696365;
    let note_root = 0x2f36d4404719a30512af45be47c9732e916cb131933102b04ba6432602db209c;
    let index = 0;
    let note_hash_path = [
        0x1e61bdae0f027b1b2159e1f9d3f8d00fa668a952dddd822fda80dc745d6f65cc,
        0x0e4223f3925f98934393c74975142bd73079ab0621f4ee133cee050a3c194f1a,
        0x2fd7bb412155bf8693a3bd2a3e7581a679c95c68a052f835dddca85fa1569a40
    ];
    let secret = 0x1929ea3ab8d9106a899386883d9428f8256cfedb3c4f6b66bf4aa4d28a79988f;
    let nullifierHash = 0x2b5872c4bee7d321717e3849b1093039b7b2da26065eb9039465fc632b14562c;

    let res = main(
        recipient,
        priv_key,
        note_root,
        index,
        note_hash_path,
        secret,
        nullifierHash
    );
    println(f"Shielded output: {res}");
}