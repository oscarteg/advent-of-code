use std::collections::HashMap;

use anyhow::{bail, ensure, Context, Result};

/// The shape to look for once the image is assembled.
const SEA_MONSTER: [&str; 3] = [
    "..................#.",
    "#....##....##....###",
    ".#..#..#..#..#..#...",
];

/// One side of an image, read left to right or top to bottom.
type Edge = Vec<bool>;

/// Canonical edge -> indices of the tiles owning it. An edge owned by a single
/// tile cannot be joined to anything and so lies on the border of the image.
type EdgeOwners = HashMap<Edge, Vec<usize>>;

/// A square bitmap of pixels, set where the input has a hash. Never empty.
#[derive(Debug, Clone, PartialEq, Eq)]
struct Image {
    rows: Vec<Vec<bool>>,
}

impl Image {
    fn parse(lines: &[&str]) -> Result<Self> {
        let parse_line = |line: &&str| -> Result<Vec<bool>> {
            line.chars()
                .map(|pixel| match pixel {
                    '#' => Ok(true),
                    '.' => Ok(false),
                    other => bail!("invalid pixel {other:?}"),
                })
                .collect()
        };
        let rows: Vec<Vec<bool>> = lines.iter().map(parse_line).collect::<Result<_>>()?;
        let size = rows.len();
        ensure!(size >= 3, "a tile needs at least 3 rows, got {size}");
        for row in &rows {
            ensure!(
                row.len() == size,
                "tile is not square: {size} rows but a row of {}",
                row.len()
            );
        }
        Ok(Self { rows })
    }

    fn size(&self) -> usize {
        self.rows.len()
    }

    fn top(&self) -> Edge {
        self.rows[0].clone()
    }

    fn bottom(&self) -> Edge {
        self.rows[self.size() - 1].clone()
    }

    fn left(&self) -> Edge {
        self.rows.iter().map(|row| row[0]).collect()
    }

    fn right(&self) -> Edge {
        self.rows.iter().map(|row| row[row.len() - 1]).collect()
    }

    fn edges(&self) -> [Edge; 4] {
        [self.top(), self.bottom(), self.left(), self.right()]
    }

    /// A quarter turn clockwise.
    fn rotated(&self) -> Self {
        let size = self.size();
        Self {
            rows: (0..size)
                .map(|row| {
                    (0..size)
                        .map(|column| self.rows[size - 1 - column][row])
                        .collect()
                })
                .collect(),
        }
    }

    /// Mirrored left to right.
    fn flipped(&self) -> Self {
        Self {
            rows: self
                .rows
                .iter()
                .map(|row| row.iter().rev().copied().collect())
                .collect(),
        }
    }

    /// The four rotations of this image followed by the four of its mirror image.
    fn orientations(&self) -> Vec<Self> {
        let mut orientations = Vec::with_capacity(8);
        for start in [self.clone(), self.flipped()] {
            let mut current = start;
            for _ in 0..4 {
                orientations.push(current.clone());
                current = current.rotated();
            }
        }
        orientations
    }

    /// This image without its one pixel wide border.
    fn interior(&self) -> Self {
        let last = self.size() - 1;
        Self {
            rows: self.rows[1..last]
                .iter()
                .map(|row| row[1..last].to_vec())
                .collect(),
        }
    }

    fn lit_pixels(&self) -> usize {
        self.rows.iter().flatten().filter(|&&pixel| pixel).count()
    }

    /// How often the pattern sits entirely on lit pixels of this image.
    fn occurrences(&self, pattern: &Pattern) -> usize {
        let size = self.size();
        if size < pattern.height || size < pattern.width {
            return 0;
        }
        (0..=size - pattern.height)
            .flat_map(|row| (0..=size - pattern.width).map(move |column| (row, column)))
            .filter(|&(row, column)| {
                pattern
                    .lit
                    .iter()
                    .all(|&(down, right)| self.rows[row + down][column + right])
            })
            .count()
    }
}

/// A shape that need not be square, stored as the offsets of its lit pixels.
struct Pattern {
    lit: Vec<(usize, usize)>,
    height: usize,
    width: usize,
}

impl Pattern {
    fn new(lines: &[&str]) -> Self {
        Self {
            lit: lines
                .iter()
                .enumerate()
                .flat_map(|(row, line)| {
                    line.char_indices()
                        .filter(|&(_, pixel)| pixel == '#')
                        .map(move |(column, _)| (row, column))
                })
                .collect(),
            height: lines.len(),
            width: lines.iter().map(|line| line.len()).max().unwrap_or(0),
        }
    }
}

#[derive(Debug, Clone)]
struct Tile {
    id: u64,
    image: Image,
}

impl Tile {
    fn parse(block: &str) -> Result<Self> {
        let mut lines = block.lines().map(str::trim).filter(|line| !line.is_empty());
        let header = lines.next().context("empty tile block")?;
        let id = header
            .strip_prefix("Tile ")
            .and_then(|rest| rest.strip_suffix(':'))
            .with_context(|| format!("expected a \"Tile <id>:\" header, got {header:?}"))?
            .parse::<u64>()
            .with_context(|| format!("invalid tile id in {header:?}"))?;
        let lines: Vec<&str> = lines.collect();
        let image = Image::parse(&lines).with_context(|| format!("tile {id}"))?;
        Ok(Self { id, image })
    }
}

fn parse_tiles(input: &str) -> Result<Vec<Tile>> {
    let tiles: Vec<Tile> = input
        .split("\n\n")
        .filter(|block| !block.trim().is_empty())
        .map(Tile::parse)
        .collect::<Result<_>>()?;
    ensure!(!tiles.is_empty(), "no tiles in input");
    let size = tiles[0].image.size();
    ensure!(
        tiles.iter().all(|tile| tile.image.size() == size),
        "tiles differ in size"
    );
    let side = tiles.len().isqrt();
    ensure!(
        side * side == tiles.len(),
        "{} tiles do not form a square image",
        tiles.len()
    );
    Ok(tiles)
}

/// The same value for both reading directions of an edge.
fn canonical(edge: &[bool]) -> Edge {
    let backward: Edge = edge.iter().rev().copied().collect();
    let forward = edge.to_vec();
    forward.min(backward)
}

fn edge_owners(tiles: &[Tile]) -> EdgeOwners {
    let mut owners = EdgeOwners::new();
    for (index, tile) in tiles.iter().enumerate() {
        for edge in tile.image.edges() {
            owners.entry(canonical(&edge)).or_default().push(index);
        }
    }
    owners
}

fn is_border(owners: &EdgeOwners, edge: &[bool]) -> bool {
    owners.get(&canonical(edge)).map_or(0, Vec::len) == 1
}

/// Corners are the tiles with exactly two edges no other tile can be joined to.
fn corner_indices(tiles: &[Tile], owners: &EdgeOwners) -> Result<Vec<usize>> {
    let corners: Vec<usize> = tiles
        .iter()
        .enumerate()
        .filter(|(_, tile)| {
            tile.image
                .edges()
                .into_iter()
                .filter(|edge| is_border(owners, edge))
                .count()
                == 2
        })
        .map(|(index, _)| index)
        .collect();
    ensure!(
        corners.len() == 4,
        "expected 4 corner tiles, found {}",
        corners.len()
    );
    Ok(corners)
}

/// One position of the layout and the edges a tile must present there.
struct Slot<'a> {
    row: usize,
    column: usize,
    side: usize,
    above: Option<&'a [bool]>,
    left: Option<&'a [bool]>,
}

impl Slot<'_> {
    /// An edge shared with a placed neighbour, owned by just one other tile.
    fn joined_edge(&self) -> Option<&[bool]> {
        self.left.or(self.above)
    }

    fn accepts(&self, image: &Image, owners: &EdgeOwners) -> bool {
        let joins = |edge: Edge, neighbour: Option<&[bool]>| match neighbour {
            Some(expected) => edge.as_slice() == expected,
            None => is_border(owners, &edge),
        };
        joins(image.top(), self.above)
            && joins(image.left(), self.left)
            && (self.row + 1 < self.side || is_border(owners, &image.bottom()))
            && (self.column + 1 < self.side || is_border(owners, &image.right()))
    }
}

/// Lays the tiles out row by row and returns the stitched image.
fn assemble(tiles: &[Tile], owners: &EdgeOwners, corners: &[usize]) -> Result<Image> {
    let side = tiles.len().isqrt();
    let mut layout: Vec<Image> = Vec::with_capacity(tiles.len());
    let mut used = vec![false; tiles.len()];

    for position in 0..tiles.len() {
        let (row, column) = (position / side, position % side);
        let above = (row > 0).then(|| layout[position - side].bottom());
        let left = (column > 0).then(|| layout[position - 1].right());
        let slot = Slot {
            row,
            column,
            side,
            above: above.as_deref(),
            left: left.as_deref(),
        };
        let candidates = match slot.joined_edge() {
            Some(edge) => owners.get(&canonical(edge)).cloned().unwrap_or_default(),
            None => corners.to_vec(),
        };
        let mut fitting = candidates
            .into_iter()
            .filter(|&index| !used[index])
            .flat_map(|index| {
                tiles[index]
                    .image
                    .orientations()
                    .into_iter()
                    .map(move |image| (index, image))
            })
            .filter(|(_, image)| slot.accepts(image, owners))
            .collect::<Vec<_>>()
            .into_iter();
        let (index, image) = fitting
            .next()
            .with_context(|| format!("no tile fits at row {row}, column {column}"))?;
        // Every slot but the first is pinned by a neighbour edge, so a second fit there
        // means the tiles can be laid out in more than one way. The first slot accepts
        // any corner in either of its two border-aligned orientations; all eight choices
        // assemble the same image up to rotation and reflection, which both parts ignore.
        if position > 0 {
            ensure!(
                fitting.next().is_none(),
                "more than one tile fits at row {row}, column {column}"
            );
        }
        used[index] = true;
        layout.push(image);
    }

    Ok(stitch(&layout, side))
}

/// Joins the laid out tiles into one image, dropping the borders that lined them up.
fn stitch(layout: &[Image], side: usize) -> Image {
    let mut rows: Vec<Vec<bool>> = Vec::new();
    for band in layout.chunks(side) {
        let interiors: Vec<Image> = band.iter().map(Image::interior).collect();
        let height = interiors.first().map_or(0, Image::size);
        for row in 0..height {
            rows.push(
                interiors
                    .iter()
                    .flat_map(|interior| interior.rows[row].iter().copied())
                    .collect(),
            );
        }
    }
    Image { rows }
}

/// Lit pixels left once every sea monster is marked. Monsters appear in exactly one
/// orientation of the image, so the orientation with the most of them is the right one.
fn water_roughness(image: &Image) -> usize {
    let monster = Pattern::new(&SEA_MONSTER);
    let monsters = image
        .orientations()
        .iter()
        .map(|view| view.occurrences(&monster))
        .max()
        .unwrap_or(0);
    image
        .lit_pixels()
        .saturating_sub(monsters * monster.lit.len())
}

pub fn solve(input: &str) -> Result<(u64, usize)> {
    let tiles = parse_tiles(input)?;
    let owners = edge_owners(&tiles);
    let corners = corner_indices(&tiles, &owners)?;
    let corner_product: u64 = corners.iter().map(|&index| tiles[index].id).product();
    let image = assemble(&tiles, &owners, &corners)?;
    Ok((corner_product, water_roughness(&image)))
}

#[cfg(test)]
mod tests {
    use super::*;

    const EXAMPLE: &str = include_str!("../../input/2020/day20-example.txt");
    const EXAMPLE_ANSWERS: (u64, usize) = (20_899_048_083_289, 273);

    /// Mirrors every tile of a jigsaw; tiles arrive in a random orientation anyway.
    fn flip_tiles(input: &str) -> String {
        input
            .lines()
            .map(|line| {
                if line.starts_with("Tile ") {
                    line.to_string()
                } else {
                    line.chars().rev().collect()
                }
            })
            .collect::<Vec<String>>()
            .join("\n")
    }

    #[test]
    fn solves_the_example() {
        assert_eq!(solve(EXAMPLE).unwrap(), EXAMPLE_ANSWERS);
    }

    #[test]
    fn answers_ignore_tile_order_and_orientation() {
        let mut blocks: Vec<String> = EXAMPLE.trim().split("\n\n").map(flip_tiles).collect();
        blocks.reverse();
        assert_eq!(solve(&blocks.join("\n\n")).unwrap(), EXAMPLE_ANSWERS);
    }

    #[test]
    fn four_rotations_return_the_original_image() {
        let image = Image::parse(&["#..", ".#.", "..#"]).unwrap();
        assert_eq!(image.rotated().rotated().rotated().rotated(), image);
        assert_eq!(image.rotated().top(), vec![false, false, true]);
        assert_eq!(image.flipped().top(), vec![false, false, true]);
        assert_eq!(image.interior().rows, vec![vec![true]]);
    }

    #[test]
    fn finds_the_corners_of_the_example() {
        let tiles = parse_tiles(EXAMPLE).unwrap();
        let owners = edge_owners(&tiles);
        let mut ids: Vec<u64> = corner_indices(&tiles, &owners)
            .unwrap()
            .iter()
            .map(|&index| tiles[index].id)
            .collect();
        ids.sort_unstable();
        assert_eq!(ids, vec![1171, 1951, 2971, 3079]);
    }

    #[test]
    fn rejects_empty_input() {
        assert!(solve("   ").is_err());
    }

    #[test]
    fn rejects_a_missing_tile_header() {
        assert!(solve("..##.#..#.\n##..#.....\n#...##..#.\n").is_err());
    }

    #[test]
    fn rejects_an_unparsable_tile_id() {
        assert!(solve(&EXAMPLE.replacen("Tile 2311:", "Tile two:", 1)).is_err());
    }

    #[test]
    fn rejects_an_unknown_pixel() {
        assert!(solve(&EXAMPLE.replacen('.', "x", 1)).is_err());
    }

    #[test]
    fn rejects_a_ragged_tile() {
        assert!(solve(&EXAMPLE.replacen("..##.#..#.\n", "..##.#..#\n", 1)).is_err());
    }

    #[test]
    fn rejects_an_incomplete_image() {
        let blocks: Vec<&str> = EXAMPLE.trim().split("\n\n").take(8).collect();
        assert!(solve(&blocks.join("\n\n")).is_err());
    }
}
