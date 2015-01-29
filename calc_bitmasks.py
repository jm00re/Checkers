# Remove left most black squares
print hex(0b11101111111011111110111111101111)
# Remove right most black squares
print hex(0b11110111111101111111011111110111)

# Remove two columns of right most black squares
print hex(0b01110111011101110111011101110111)
# Remove two columns of left most black squares
print hex(0b11101110111011101110111011101110)

# Remove first 4 black squares
print hex(0b00001111111111111111111111111111)
# Remove last 4 black squares
print hex(0b11111111111111111111111111110000)

# Remove all but first 4 black squares
print hex(0b1111)
# Remove all but law 4 black squares
print hex(0b11110000000000000000000000000000)

# Even rows
print hex(0b00001111000011110000111100001111)
# Off rows
print hex(0b11110000111100001111000011110000)
