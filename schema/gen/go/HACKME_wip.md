
### absent values

The handling of absent values is still not consistent.

Currently:

- reading (via accessors or iterators) yields `ipld.Undef` values for absent fields
- putting those ipld.Undef values via NodeAssembler.AssignNode will result in `ErrWrongKind`.
- *the recursive copies embedded in AssignNode methods don't handle undefs either_.

The first two are defensible and consistent (if not necessarily ergonomic).
The third is downright a bug, and needs to be fixed.

How we fix it is not entirely settled.

- Option 1: keep the hostility to undef assignment
- Option 2: *require* explicit undef assignment
- Option 3: become indifferent to undef assignment when it's valid
- Option 4: don't yield values that are undef during iteration at all

Option 3 seems the most preferrable (and least user-hostile).
(Options 1 and 2 create work for end users;
Option 4 has questionable logical consistency.)

Updating the codegen to do Option 3 needs some work, though.

It's likely that the way to go about this would involve adding two more valid
bit states to the extended schema.Maybe values: one for allowAbsent (similar to
the existing allowNull), and another for both (for "nullable optional" fields).
Every NodeAssembler would then have to support that, just as they each support allowNull now.

I think the above design is valid, but it's not implemented nor tested yet.
