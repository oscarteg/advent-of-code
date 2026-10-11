use anyhow::{anyhow, ensure, Context, Result};
use std::collections::{HashMap, HashSet};

/// Beacons two scanners must share before their reports can be merged.
const MIN_OVERLAP: usize = 12;

type Point = [i32; 3];

/// The beacons one scanner sees, in that scanner's own frame of reference.
#[derive(Debug, Default)]
struct Report {
    beacons: Vec<Point>,
}

/// A signed axis permutation with determinant +1: one scanner orientation.
#[derive(Clone, Copy)]
struct Rotation {
    axes: [usize; 3],
    signs: [i32; 3],
}

impl Rotation {
    /// The 24 orientations: six axis permutations combined with the sign flips that
    /// keep the determinant positive, which drops the 24 mirror images.
    fn all() -> Vec<Self> {
        const PERMUTATIONS: [([usize; 3], i32); 6] = [
            ([0, 1, 2], 1),
            ([1, 2, 0], 1),
            ([2, 0, 1], 1),
            ([0, 2, 1], -1),
            ([2, 1, 0], -1),
            ([1, 0, 2], -1),
        ];
        let mut rotations = Vec::with_capacity(24);
        for (axes, parity) in PERMUTATIONS {
            for flips in 0..8 {
                let signs = [0, 1, 2].map(|axis| if (flips & (1 << axis)) == 0 { 1 } else { -1 });
                if parity * signs.iter().product::<i32>() == 1 {
                    rotations.push(Self { axes, signs });
                }
            }
        }
        rotations
    }

    fn apply(self, point: Point) -> Point {
        [0, 1, 2].map(|axis| self.signs[axis] * point[self.axes[axis]])
    }
}

fn parse_point(line: &str) -> Result<Point> {
    let coordinates = line
        .split(',')
        .map(|value| {
            value
                .trim()
                .parse::<i32>()
                .with_context(|| format!("invalid coordinate {value:?} in {line:?}"))
        })
        .collect::<Result<Vec<_>>>()?;
    coordinates.try_into().map_err(|values: Vec<i32>| {
        anyhow!(
            "beacon {line:?} needs three coordinates, got {}",
            values.len()
        )
    })
}

fn parse(input: &str) -> Result<Vec<Report>> {
    let mut reports: Vec<Report> = Vec::new();
    for line in input.lines().map(str::trim).filter(|line| !line.is_empty()) {
        if let Some(id) = line
            .strip_prefix("--- scanner ")
            .and_then(|rest| rest.strip_suffix(" ---"))
        {
            id.parse::<u32>()
                .with_context(|| format!("invalid scanner id in header {line:?}"))?;
            reports.push(Report::default());
            continue;
        }
        let report = reports
            .last_mut()
            .with_context(|| format!("beacon {line:?} appears before any scanner header"))?;
        report.beacons.push(parse_point(line)?);
    }
    Ok(reports)
}

fn translate(point: Point, offset: Point) -> Point {
    [0, 1, 2].map(|axis| point[axis] + offset[axis])
}

/// Searches for the orientation and offset that move a report into the frame of the
/// already placed beacons, returning the scanner position and its placed beacons.
fn align(
    beacons: &HashSet<Point>,
    report: &Report,
    rotations: &[Rotation],
) -> Option<(Point, Vec<Point>)> {
    for rotation in rotations {
        let rotated: Vec<Point> = report
            .beacons
            .iter()
            .map(|&beacon| rotation.apply(beacon))
            .collect();
        let mut offsets: HashMap<Point, usize> = HashMap::new();
        for &known in beacons {
            for &beacon in &rotated {
                let offset = [0, 1, 2].map(|axis| known[axis] - beacon[axis]);
                let matches = offsets.entry(offset).or_default();
                *matches += 1;
                if *matches >= MIN_OVERLAP {
                    let placed = rotated.iter().map(|&b| translate(b, offset)).collect();
                    return Some((offset, placed));
                }
            }
        }
    }
    None
}

/// Places every report in scanner 0's frame, repeatedly matching the reports that are
/// still unplaced against the beacons of a freshly placed one.
fn assemble(reports: &[Report]) -> Result<(HashSet<Point>, Vec<Point>)> {
    let (first, rest) = reports
        .split_first()
        .context("input contains no scanner reports")?;
    let rotations = Rotation::all();
    let mut beacons: HashSet<Point> = first.beacons.iter().copied().collect();
    let mut scanners = vec![[0, 0, 0]];
    let mut frontier = vec![first.beacons.clone()];
    let mut pending: Vec<&Report> = rest.iter().collect();

    while let Some(placed) = frontier.pop() {
        let placed: HashSet<Point> = placed.into_iter().collect();
        let mut unplaced = Vec::with_capacity(pending.len());
        for report in std::mem::take(&mut pending) {
            match align(&placed, report, &rotations) {
                Some((scanner, aligned)) => {
                    beacons.extend(aligned.iter().copied());
                    scanners.push(scanner);
                    frontier.push(aligned);
                }
                None => unplaced.push(report),
            }
        }
        pending = unplaced;
    }

    ensure!(
        pending.is_empty(),
        "{} scanner reports share fewer than {MIN_OVERLAP} beacons with the rest",
        pending.len()
    );
    Ok((beacons, scanners))
}

fn manhattan(from: Point, to: Point) -> i32 {
    (0..3).map(|axis| (from[axis] - to[axis]).abs()).sum()
}

fn largest_distance(scanners: &[Point]) -> i32 {
    scanners
        .iter()
        .enumerate()
        .flat_map(|(index, &from)| {
            scanners[index + 1..]
                .iter()
                .map(move |&to| manhattan(from, to))
        })
        .max()
        .unwrap_or(0)
}

pub fn solve(input: &str) -> Result<(usize, i32)> {
    let reports = parse(input)?;
    let (beacons, scanners) = assemble(&reports)?;
    Ok((beacons.len(), largest_distance(&scanners)))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/day19-example.txt");

    /// The twelve beacons scanners 0 and 1 share, as the puzzle lists them in the frame
    /// of each scanner; scanner 1 sits at 68,-1246,-43.
    const SHARED_TWELVE: &str = "--- scanner 0 ---
-618,-824,-621
-537,-823,-458
-447,-329,318
404,-588,-901
544,-627,-890
528,-643,409
-661,-816,-575
390,-675,-793
423,-701,434
-345,-311,381
459,-707,401
-485,-357,347

--- scanner 1 ---
686,422,578
605,423,415
515,917,-361
-336,658,858
-476,619,847
-460,603,-452
729,430,532
-322,571,750
-355,545,-477
413,935,-424
-391,539,-444
553,889,-390
";

    /// The puzzle's six beacons written in five of the 24 orientations.
    const ORIENTATIONS: &str = "--- scanner 0 ---
-1,-1,1
-2,-2,2
-3,-3,3
-2,-3,1
5,6,-4
8,0,7

--- scanner 0 ---
1,-1,1
2,-2,2
3,-3,3
2,-1,3
-5,4,-6
-8,-7,0

--- scanner 0 ---
-1,-1,-1
-2,-2,-2
-3,-3,-3
-1,-3,-2
4,6,5
-7,0,8

--- scanner 0 ---
1,1,-1
2,2,-2
3,3,-3
1,3,-2
-4,-6,5
7,0,8

--- scanner 0 ---
1,1,1
2,2,2
3,3,3
3,1,2
-6,-4,-5
0,7,-8
";

    #[test]
    fn solves_official_example() {
        assert_eq!(solve(EXAMPLE).unwrap(), (79, 3621));
    }

    #[test]
    fn merges_a_pair_sharing_exactly_twelve_beacons() {
        assert_eq!(solve(SHARED_TWELVE).unwrap(), (12, 1357));
    }

    #[test]
    fn single_report_needs_no_alignment() {
        assert_eq!(solve("--- scanner 0 ---\n1,2,3\n4,5,6\n").unwrap(), (2, 0));
    }

    #[test]
    fn rotations_reorient_one_beacon_set_onto_itself() {
        let rotations = Rotation::all();
        assert_eq!(rotations.len(), 24);
        let reports = parse(ORIENTATIONS).unwrap();
        let reference: HashSet<Point> = reports[0].beacons.iter().copied().collect();
        for report in &reports[1..] {
            assert!(rotations.iter().any(|rotation| {
                report
                    .beacons
                    .iter()
                    .map(|&beacon| rotation.apply(beacon))
                    .collect::<HashSet<_>>()
                    == reference
            }));
        }
    }

    #[test]
    fn rejects_malformed_input() {
        for input in [
            "",
            "1,2,3\n",
            "--- scanner 0 ---\n1,2\n",
            "--- scanner 0 ---\n1,2,x\n",
            "--- scanner one ---\n1,2,3\n",
            "--- scanner 0 ---\n1,2,3\n\n--- scanner 1 ---\n4,5,6\n",
        ] {
            assert!(solve(input).is_err(), "expected an error for {input:?}");
        }
    }
}
