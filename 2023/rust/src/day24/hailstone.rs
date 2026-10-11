use crate::day24::{
    custom_error::AocError,
    ratio::{gcd, Ratio},
};

/// A position or velocity with exact integer coordinates.
#[derive(Clone, Copy, Debug)]
pub struct Vector {
    pub x: i128,
    pub y: i128,
    pub z: i128,
}

impl Vector {
    const fn components(self) -> [i128; 3] {
        [self.x, self.y, self.z]
    }

    fn checked_sub(self, other: Self) -> Result<Self, AocError> {
        let difference =
            |left: i128, right: i128| left.checked_sub(right).ok_or(AocError::Overflow);
        Ok(Self {
            x: difference(self.x, other.x)?,
            y: difference(self.y, other.y)?,
            z: difference(self.z, other.z)?,
        })
    }

    fn cross(self, other: Self) -> Result<Self, AocError> {
        Ok(Self {
            x: determinant(self.y, self.z, other.y, other.z)?,
            y: determinant(self.z, self.x, other.z, other.x)?,
            z: determinant(self.x, self.y, other.x, other.y)?,
        })
    }

    /// Divides the components by their greatest common divisor, which keeps the
    /// line the vector spans while shrinking the numbers. `None` for the zero
    /// vector, which spans no line.
    fn primitive(self) -> Option<Self> {
        let divisor = self.components().into_iter().fold(0, gcd);
        Some(Self {
            x: self.x.checked_div(divisor)?,
            y: self.y.checked_div(divisor)?,
            z: self.z.checked_div(divisor)?,
        })
    }

    fn coordinate_sum(self) -> Result<i128, AocError> {
        self.x
            .checked_add(self.y)
            .and_then(|partial| partial.checked_add(self.z))
            .ok_or(AocError::Overflow)
    }
}

/// The determinant `a * d - b * c` of the matrix `[[a, b], [c, d]]`.
fn determinant(a: i128, b: i128, c: i128, d: i128) -> Result<i128, AocError> {
    a.checked_mul(d)
        .zip(b.checked_mul(c))
        .and_then(|(left, right)| left.checked_sub(right))
        .ok_or(AocError::Overflow)
}

/// A hailstone's position at time zero and its constant velocity.
#[derive(Clone, Copy, Debug)]
pub struct Hailstone {
    pub position: Vector,
    pub velocity: Vector,
}

/// Parses one hailstone per line, each written as `px, py, pz @ vx, vy, vz`.
pub fn parse(input: &str) -> Result<Vec<Hailstone>, AocError> {
    let stones = input
        .lines()
        .map(str::trim)
        .filter(|line| !line.is_empty())
        .map(parse_hailstone)
        .collect::<Result<Vec<_>, _>>()?;
    if stones.is_empty() {
        return Err(AocError::InvalidInput("no hailstones"));
    }
    Ok(stones)
}

fn parse_hailstone(line: &str) -> Result<Hailstone, AocError> {
    let (position, velocity) = line
        .split_once('@')
        .ok_or(AocError::InvalidInput("a hailstone needs an '@' separator"))?;
    Ok(Hailstone {
        position: parse_vector(position)?,
        velocity: parse_vector(velocity)?,
    })
}

fn parse_vector(input: &str) -> Result<Vector, AocError> {
    let coordinates = input
        .split(',')
        .map(|coordinate| coordinate.trim().parse::<i64>())
        .collect::<Result<Vec<_>, _>>()?;
    let [x, y, z] = coordinates.as_slice() else {
        return Err(AocError::InvalidInput("a vector needs three coordinates"));
    };
    Ok(Vector {
        x: i128::from(*x),
        y: i128::from(*y),
        z: i128::from(*z),
    })
}

/// Counts the pairs of hailstones whose forward paths cross inside the square
/// `[low, high]` on both axes, ignoring the z axis.
pub fn crossings_in_area(stones: &[Hailstone], low: i128, high: i128) -> Result<usize, AocError> {
    let mut crossings = 0;
    for (index, first) in stones.iter().enumerate() {
        for second in &stones[index + 1..] {
            if crosses_in_area(*first, *second, low, high)? {
                crossings += 1;
            }
        }
    }
    Ok(crossings)
}

fn crosses_in_area(
    first: Hailstone,
    second: Hailstone,
    low: i128,
    high: i128,
) -> Result<bool, AocError> {
    // Cramer's rule on `first.position + s * first.velocity == second.position
    // + t * second.velocity`, restricted to the x and y axes.
    let denominator = determinant(
        first.velocity.x,
        first.velocity.y,
        second.velocity.x,
        second.velocity.y,
    )?;
    if denominator == 0 {
        return Ok(false);
    }
    let offset = second.position.checked_sub(first.position)?;
    let first_time = determinant(offset.x, offset.y, second.velocity.x, second.velocity.y)?;
    let second_time = determinant(offset.x, offset.y, first.velocity.x, first.velocity.y)?;
    // Both times and the crossing itself stay scaled by the denominator so the
    // comparisons need no division; a negative denominator would reverse them.
    let (denominator, first_time, second_time) = if denominator < 0 {
        (
            denominator.checked_neg().ok_or(AocError::Overflow)?,
            first_time.checked_neg().ok_or(AocError::Overflow)?,
            second_time.checked_neg().ok_or(AocError::Overflow)?,
        )
    } else {
        (denominator, first_time, second_time)
    };
    if first_time <= 0 || second_time <= 0 {
        return Ok(false);
    }
    let bounds = low.checked_mul(denominator).ok_or(AocError::Overflow)?
        ..=high.checked_mul(denominator).ok_or(AocError::Overflow)?;
    for (position, velocity) in [
        (first.position.x, first.velocity.x),
        (first.position.y, first.velocity.y),
    ] {
        let coordinate = position
            .checked_mul(denominator)
            .zip(first_time.checked_mul(velocity))
            .and_then(|(start, travelled)| start.checked_add(travelled))
            .ok_or(AocError::Overflow)?;
        if !bounds.contains(&coordinate) {
            return Ok(false);
        }
    }
    Ok(true)
}

/// A hailstone's path in the first hailstone's frame of reference, where the
/// rock's own path runs through the origin.
struct RelativePath {
    offset: Vector,
    velocity: Vector,
    /// Normal of the plane spanned by this path and the origin.
    normal: Vector,
}

impl RelativePath {
    fn position_at(&self, time: Ratio, axis: usize) -> Result<Ratio, AocError> {
        let offset = Ratio::from_integer(self.offset.components()[axis]);
        let velocity = Ratio::from_integer(self.velocity.components()[axis]);
        offset.checked_add(time.checked_mul(velocity)?)
    }
}

/// Sums the coordinates of the position the rock has to be thrown from.
///
/// Seen from the first hailstone the rock's path runs through the origin, so
/// every other hailstone's path spans a plane with it. The rock's direction has
/// to lie in all of those planes, which makes it the cross product of two of
/// their normals. Crossing that line with two hailstone paths gives two hit
/// times and two points, and walking back from them to time zero yields the
/// rock's velocity and start position.
pub fn rock_position_sum(stones: &[Hailstone]) -> Result<i128, AocError> {
    let Some((base, rest)) = stones.split_first() else {
        return Err(AocError::InvalidInput("no hailstones"));
    };
    let paths = relative_paths(*base, rest)?;
    let (first, second, direction) = rock_direction(&paths)?;
    let first_time = hit_time(first, direction)?;
    let second_time = hit_time(second, direction)?;
    let span = second_time.checked_sub(first_time)?;
    if span.is_zero() {
        return Err(AocError::InvalidInput(
            "two hailstones are hit at the same time",
        ));
    }
    let mut total = Ratio::from_integer(base.position.coordinate_sum()?);
    for axis in 0..3 {
        let first_point = first.position_at(first_time, axis)?;
        let second_point = second.position_at(second_time, axis)?;
        let velocity = second_point.checked_sub(first_point)?.checked_div(span)?;
        total = total.checked_add(first_point.checked_sub(first_time.checked_mul(velocity)?)?)?;
    }
    total.to_integer().ok_or(AocError::InvalidInput(
        "the rock does not start on integer coordinates",
    ))
}

fn relative_paths(base: Hailstone, rest: &[Hailstone]) -> Result<Vec<RelativePath>, AocError> {
    let mut paths = Vec::with_capacity(rest.len());
    for stone in rest {
        let offset = stone.position.checked_sub(base.position)?;
        let velocity = stone.velocity.checked_sub(base.velocity)?;
        // A path through the origin spans no plane with it, so it says nothing
        // about the rock's direction.
        if let Some(normal) = offset.cross(velocity)?.primitive() {
            paths.push(RelativePath {
                offset,
                velocity,
                normal,
            });
        }
    }
    Ok(paths)
}

/// Picks two paths whose planes meet in the rock's line and returns that line's
/// direction alongside them.
fn rock_direction(
    paths: &[RelativePath],
) -> Result<(&RelativePath, &RelativePath, Vector), AocError> {
    for (index, first) in paths.iter().enumerate() {
        for second in &paths[index + 1..] {
            if let Some(direction) = first.normal.cross(second.normal)?.primitive() {
                return Ok((first, second, direction));
            }
        }
    }
    Err(AocError::InvalidInput(
        "the hailstones do not determine a single rock",
    ))
}

/// When the path crosses the rock's line: the point `offset + time * velocity`
/// lies on that line exactly when its cross product with the direction is zero.
fn hit_time(path: &RelativePath, direction: Vector) -> Result<Ratio, AocError> {
    let offset = path.offset.cross(direction)?.components();
    let velocity = path.velocity.cross(direction)?.components();
    let axis = velocity
        .iter()
        .position(|&component| component != 0)
        .ok_or(AocError::InvalidInput(
            "a hailstone runs parallel to the rock",
        ))?;
    Ratio::new(
        offset[axis].checked_neg().ok_or(AocError::Overflow)?,
        velocity[axis],
    )
}
