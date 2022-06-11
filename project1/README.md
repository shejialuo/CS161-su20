# Project 1

This project will use [32-bit virtual machine](https://drive.google.com/file/d/1Rni4mp7nEdh8-jok1x-__p4ADVRNNdmy).

## Task 1

Task 1 uses stack smashing to do code injection. In order to find the byte
we need to write. User `gdb` to look at the start address of `door`.

![The result of the gdb](./assets/The%20result%20of%20the%20gdb.png)

## Task 2

Task2 uses tack smashing to do code injection. The idea is that `fread` except
`unsigned`, however the type of `size` is `int8_t`, so we could make the `size`
to be negative, thus `fread` would convert `size` to unsigned to make it a big
number thus we can inject whatever we want.

## Task 3

Task3 adds the stack canary, so in `interact` we need to first find the stack
canary, which is the most difficult question. However, we could exploit the
following code snippet:

```c
gets(c.buffer)
while (c.buffer[i]) {
  if (c.buffer[i] == '\\' && c.buffer[i+1] == 'x') {
    int top_half = nibble_to_int(c.buffer[i+2]);
    int bottom_half = nibble_to_int(c.buffer[i+3]);
    c.answer[j] = top_half << 4 | bottom_half;
    i += 3;
  }
    else {
      c.answer[j] = c.buffer[i];
    }
  i++; j++;
}
c.answer[j] = 0;
printf("%s\n", c.answer);
```

In this code snippet, the thing we need to utilize is the `printf("%s\n", c.answer)`.
The stack canary is 4 bytes with deterministic 00 one byte, so we should cross the
null byte of the stack canary, thus `c.answer` would print the stack canary we want.

What we could utilize is the `i += 3`. So we could just pass `\x` to the input at
the last, which makes it cross the null byte of the stack canary. So we could find
stack canary in `c.answer`.
