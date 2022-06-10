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
