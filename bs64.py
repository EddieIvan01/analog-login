class Base64(object):
    def __init__(self):
        self.b64map = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
        self.b64pad = "="

    
    def integer_to_char(self,integer):

        if 0 <= integer <= 9:
            return str(integer)
        elif 9 < integer <= 35:
            return chr(integer+87)
        else:
            return ""

    def hex_to_b64(self, string):
        i = 0
        c = 0
        ret = ""
        while i + 3 <= len(string):
            c = int(string[i:i+3], 16)
            ret += self.b64map[c >> 6] + self.b64map[c & 63]
            i = i + 3
        if i+1 == len(string):
            c = int(string[i:i+1], 16)
            ret += self.b64map[c << 2]
        elif i+2 == len(string):
            c = int(string[i:i+2], 16)
            ret += self.b64map[c >> 2] + self.b64map[(c & 3) << 4]
        while len(ret) & 3 > 0:
            ret += self.b64pad
        return ret

    def b64_to_hex(self, string):
        ret = ""
        k = 0
        slop = 0
        for i in string:
            if self.b64pad == i:
                break
            v = self.b64map.index(i)

            if v < 0:
                continue
            if k == 0:
                ret += self.integer_to_char(v >> 2)
                slop = v & 3
                k = 1
            elif k == 1:
                ret += self.integer_to_char((slop << 2) | (v >> 4))
                slop = v & 0xf
                k = 2
            elif k == 2:
                ret += self.integer_to_char(slop)
                ret += self.integer_to_char(v >> 2)
                slop = v & 3
                k = 3
            else:
                ret += self.integer_to_char((slop << 2) | (v >> 4))
                ret += self.integer_to_char(v & 0xf)
                k = 0
        if k == 1:
            ret += self.integer_to_char(slop << 2)
        return ret


if __name__ == '__main__':
    b64 = Base64()
    print(b64.b64_to_hex(
        "AJ+MzUJkH6TTpH2vCw38+2Sc5l8iUtXzG5M7+ZhLxS+kvA91AXgtJjL0VrOrvAT0EXFicUKDH/XhORfdSEmmCLJTJ79+6Il+XVX8M09Bwfw5MhONAIRf5VTTO0ryaYmhHCGaj6DAiNOC2Ob/HXyHey9Pm0/XDGRx3BQko35b/4hZ"))

    print(b64.hex_to_b64("d76df5db"))