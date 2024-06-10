use dep::std;
mod compute_root;


fn main(
    root: Field, // merkle root for the transfer
    attack_coordinates: pub [Field; 2], // coordinates of planet to attack
    missile_token_id: pub Field,
    merkle_index:  Field,
    hash_path: [Field; 3],
    priv_key: Field,
) -> pub Field {

    //prove knowledge of some priv_key and group secret
    let hashed_attack_coordinates = std::hash::pedersen_hash([ attack_coordinates[0], attack_coordinates[1] ]);
    let note_commitment = std::hash::pedersen_hash([ priv_key, hashed_attack_coordinates ]);

    // ensure the vote is member of eligible voters
    let computed_root = std::merkle::compute_merkle_root(note_commitment, merkle_index, hash_path);
    println(f"computed_root1: {computed_root}");
    println(f"root: {root}");
    assert(root == computed_root);


    // return hash to check same priv_key doesn't re-use the same missile tokenId
    let nullifier = std::hash::pedersen_hash([ root, priv_key, missile_token_id ]);
    nullifier
}



fn generate_merkle_root(priv_key: Field, index: Field, attack_coordinates: [Field; 2]) -> Field {

    let mut commitments: [Field; 8] = [0; 8];
    let hashed_attack_coordinates = std::hash::pedersen_hash([ attack_coordinates[0], attack_coordinates[1] ]);

    // level 1 leafs
    for i in 0..8 {
        if i == index as u8 {
            let commitment = std::hash::pedersen_hash([priv_key, hashed_attack_coordinates ]);
            commitments[i] = commitment;
        } else {
            let commitment = std::hash::pedersen_hash([i as Field, hashed_attack_coordinates ]);
            commitments[i] = commitment;
        }
    }

    // level 2 branches
    let branch_2a = std::hash::pedersen_hash([commitments[0], commitments[1]]);
    let branch_2b = std::hash::pedersen_hash([commitments[2], commitments[3]]);
    let branch_2c = std::hash::pedersen_hash([commitments[4], commitments[5]]);
    let branch_2d = std::hash::pedersen_hash([commitments[6], commitments[7]]);

    // level 3 branches
    let branch_3a = std::hash::pedersen_hash([branch_2a, branch_2b]);
    let branch_3b = std::hash::pedersen_hash([branch_2c, branch_2d]);

    let merkle_root = std::hash::pedersen_hash([ branch_3a, branch_3b ]);

    println("");
    println("[Merkle Tree]");
    println(f"Root:");
    println(f"{merkle_root}");
    println("");

    println(f"Level 3 Branches");
    println(f"L: {branch_3a}");
    println(f"R: {branch_3b}");
    println("");

    println(f"Level 2 Branches");
    println(f"LL: {branch_2a}");
    println(f"LR: {branch_2b}");
    println(f"RL: {branch_2c}");
    println(f"RR: {branch_2d}");
    println("");

    println(f"Leafs:");
    for i in 0..8 {
        let c = commitments[i];
        println(f"leaf {i}: {c}");
    }
    println("");
    merkle_root
}


#[test]
fn test_main() {

    // create commitment
    let attack_coordinates = [2, 3]; // coordinates of airship to attack, only player knows
    let priv_key = 0x000000000000000000000000000000000000000000000000000000616c696365; // commitment[merkle_index]
    let merkle_index =  2; // some index that server keeps track of
    let root = generate_merkle_root(priv_key, merkle_index, attack_coordinates);
    // root = 0x0324658c8f78bcb17fdcc2092a4597e939ac44c6c8889fd54a73df5e2869657d

    let missile_token_id = 1234; // use as nullifier to prevent double spending

    let hash_path = [
        // hash_paths are indexed from: bottom leafs first, left to right;
        0x299e0f032bbd01fd16326a4a33a86c6477ffc15067e799b7be4114e9a4adbc3a, // commitment_3
        0x0bd31882ca4cd0f269c3dd7ae76c0dfc397ab0a745a557de521bd693d68329aa, // LL in level 2
        0x0e258815841ece6567e5b339064066e30e2d51a2e10d8cbfbfa56aebd285d83c, // R in level 3
    ];
    //             root
    //          /        \
    //        L            R
    //     /     \       /    \
    //   LL      LR     RL     RR
    //  /  \   /  \    /  \    /  \
    // 0   1   2   3   4   5   6   7


    let nullifier_result = main(
        root,
        attack_coordinates,
        missile_token_id,
        merkle_index,
        hash_path,
        priv_key
    );
    println("");
    println(f"Missile TokenId Nullifier: {nullifier_result}");
    let expected_nullifier = 0x14c704ef63cef381567b698c7058e0e9e26b57bb92f255121b69ec58fe60186e;
    // cannot reuse that missile again, burn it.
    assert(nullifier_result == expected_nullifier);
}