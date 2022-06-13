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

## Task 4

Well, first you need to understand [off by one](http://www.icir.org/matthias/cs161-sp13/aslr-bypass.pdf). This is
a very fascinating feature. I want to conclude the idea by myself. The fault is the `strlen`
function doesn't consider about the `\0`. For example , in the `agent-brown.c`:

```c
void flip(char *buf, const char *input) {
  size_t n = strlen(input);
  int i;
  for (i = 0; i < n && i <= 64; ++i)
    buf[i] = input[i] ^ (1u << 5);

  while (i < 64)
    buf[i++] = '\0';
}

void invoke(const char *in) {
  char buf[64];
  flip(buf, in);
  puts(buf);
}

void dispatch(const char *in) {
  invoke(in);
}

```

You can see the problem, because of `i < n && i <= 64`. When the `buf` is full, the `strlen`
would return `63`, so it would still writes to next byte which does not belong to `buf`. So
the following could happen. Due to the x86 function call convention:

```assembly
enter:
  push epb
  mov epb eps

leave:
  mov eps epb
  pop epb
  pop eip
```

Well, as you can see, we can easily change the value of the top above the `buf`, to make
it to point to the `buf` inside (assumption). When the `flip` returns, nothing would be
changed. However, things are different when the `invoke` returns, because the `ebp` was
changed in the `flip`, so when the `invoke` returns, it first moves the stack pointer to
the `epb`, now `ebp` points to the `buf` content. Now we can change the `eip` to whatever
we want because we have stack pointer to the `buf`, so it pops the content to `epb` which
we do not care, but then it pops the `eip`. So we can change the value of the `eip`, thus
we can control the program.

![The memory layout](./assets/The%20memory%20layout.png)

Wonderful idea! At now, the task 4 is easy.

## Task 5

The Task 5 only examine the file size when it opens the file:

```c
void read_file() {
  char buf[MAX_BUFSIZE];
  uint32_t bytes_to_read;

  int fd = open(FILENAME, O_RDONLY);
  if (fd == -1) EXIT_WITH_ERROR("Could not find file!");

  if (file_is_too_big(fd)) EXIT_WITH_ERROR("File too big!");

  printf("How many bytes should I read? ");
  fflush(stdout);
  if (scanf("%u", &bytes_to_read) != 1)
    EXIT_WITH_ERROR("Could not read the number of bytes to read!");

  ssize_t bytes_read = read(fd, buf, bytes_to_read);
  if (bytes_read == -1) EXIT_WITH_ERROR("Could not read!");

  buf[bytes_read] = 0;
  printf("Here is the file!\n%s", buf);
  close(fd);
}
```

So we could use `opopen` system call, first we just gives a small file to make it pass
the check. We wait on the `printf("How many bytes should I read? ");`. Because the
program needs us to give the input, and this is the best time for us to change the buf
size again! Thus we could hack the system.

Well, when I have done this hack, I find that the use `ln` to create a symbol link and
use the `./dejavu` could also done:

```sh
ln -s README hack
./dejavu
```

Well, the file is too long! So ridiculous.

## Task 6

For task 6, we use `ret2esp` to hack. Because the `jump *esp` is in the `.text`,
which isn't influenced by the ASLR. It is easy to find the instruction, because
in the `magic` function, `i |= 58623`. And we can use gdb easily find the address
of the instruction is `0x8048666`. So our modified address should be `0x8048666`.

Let's we find the `buf` start address (which may be changed) is `0xbfd32570`, and
the `eip` address (which may be changed) is `0xbfd3259c`. So the padding byte is
44. Thus we could change the `eip`'s content to `0x8048666`.
