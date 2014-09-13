## Sanitize

A dead simple Go HTML whitelist-sanitization library.

### Goal

Efficiently support the following types of HTML sanitization through simple programmatic or JSON configuration:
- Removal of all non-whitelisted elements
- Unwrapping of all non-whitelisted elements

### Examples

Given a whitelist configuration

```json
{
    "div": ["id", "class"],
    "b": [],
    "i": []
}
```

and basic input

```html
<div class="my-class" style="position:relative;">
    <i>Something emphasized</i>
    <p>
        here is a
        <i>paragraph</i>
    </p>
    <b>Something bold</b> 
</div>
```

**Removal**

Removal of non-whitelisted elements in the provided example would yield

```html
<div class="my-class">
    <i>Something emphasized</i>
    <b>Something bold</b> 
</div>
```

Note how the `style` attribute was removed from the `div` element and the `p` element was removed entirely

**Unwrapping**

Unwrapping of non-whitelisted elements in the provided example would yield

```html
<div class="my-class">
    <i>Something emphasized</i>
    here is a
    <i>paragraph</i>
    <b>Something bold</b> 
</div>
```

Note how the `style` attribute was still removed from the `div` element, while the `p` element was 'unwrapped' (ie. it's children were attached to it's parent)

### Usage

### Steps to 1.0
- [x] Support sanitization that removes non-whitelisted nodes entirely
- [ ] Support sanitization that unwraps non-whitelisted nodes, allowing the text and/or whitelisted subtree through
- [ ] Whitelist-level configuration options (eg. `stripWhitespace`, ``)
- [ ] Efficient attribute checking by not allocating a new slice on every whitelisted attribute for an element
- [ ] Support non `string` type attribute values
- [ ] Tests
- [ ] Refactor configuration parsing to have io.Reader interface instead of expecting a filepath