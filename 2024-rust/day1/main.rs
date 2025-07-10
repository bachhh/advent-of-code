use std::cmp::Ordering;
use std::env;
use std::fs::File;
use std::io::{self, BufRead, BufReader};

fn main() -> io::Result<()> {
    // Get the first command-line argument (after the program name)
    let args: Vec<String> = env::args().collect();
    if args.len() < 2 {
        eprintln!("Usage: {} <filename>", args[0]);
        std::process::exit(1);
    }

    let filename = &args[1];

    // Open the file
    let file = File::open(filename)?;
    let reader = BufReader::new(file);
    let mut left: Vec<i32> = Vec::new();
    let mut right: Vec<i32> = Vec::new();

    for line in reader.lines() {
        let line = line?;
        let nums: Vec<i32> = line
            .split_whitespace()
            .take(2)
            .filter_map(|s| s.parse().ok())
            .collect();
        if nums.len() == 2 {
            left.push(nums[0]);
            right.push(nums[1]);
        }
    }
    println!("{}", part2(left, right));

    Ok(())
}

fn part1(mut left: Vec<i32>, mut right: Vec<i32>) -> i32 {
    left.sort_by(|a, b| {
        if a < b {
            Ordering::Less
        } else if a > b {
            Ordering::Greater
        } else {
            Ordering::Equal
        }
    });

    right.sort_by(|a, b| {
        if a < b {
            Ordering::Less
        } else if a > b {
            Ordering::Greater
        } else {
            Ordering::Equal
        }
    });

    let mut total = 0;
    let len = left.len().min(right.len());
    for i in 0..len {
        let diff = (left[i] - right[i]).abs();
        total += diff;
        println!("{}, {}", total, diff);
    }
    total
}

fn part2(left: Vec<i32>, right: Vec<i32>) -> i32 {
    let mut total = 0;
    for num in left.iter() {
        let mut factor = 0;
        for num2 in right.iter() {
            if num == num2 {
                factor += 1;
            }
        }
        let score = num * factor;
        println!("{}, {}, {}", num, factor, score);
        total += score;
    }
    total
}
