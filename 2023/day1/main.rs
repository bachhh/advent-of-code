use std::collections::HashMap;
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
    let mut strToNum = HashMap<string, i32>::new();

    let mut total = 0;
    // Read and print lines one by one
    for line in reader.lines() {
        let line = line?; // Handle I/O error per line

        let first: Option<char> = None;
        let mut last: Option<char> = None;
        for ch in line.chars() {
            if !ch.is_numeric() {
                continue;
            }
            if first.is_none() {
                first = Some(ch);
                last = Some(ch);
            } else {
                last = Some(ch);
            }
        }
        if let (Some(c1), Some(c2)) = (first, last) {
            if let Ok(num) = format!("{}{}", c1, c2).parse::<i32>() {
                total += num;
                println!("{}", num);
            }
        }
    }
    println!("{}", total);

    Ok(())
}
