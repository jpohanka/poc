PoC of using google cmp and cmpopts packages to create a human readable diff
for go structures.

Remaining issues:
- need to customize Reporter to output desired format
- Transformers produce wrappers that are also printed
- maybe use a different differ or write one from scratch :)
