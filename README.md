## Sanitize

A dead simple Go HTML whitelist-sanitization library.

[![wercker status](https://app.wercker.com/status/c6d103a9e8ddfa071672e62cfc2aea6a/s/master "wercker status")](https://app.wercker.com/project/bykey/c6d103a9e8ddfa071672e62cfc2aea6a)

### Goal

Efficiently support the following types of HTML sanitization through simple programmatic or JSON configuration:
- Removal of all non-whitelisted elements
- Unwrapping of all non-whitelisted elements

### Examples

Given a whitelist configuration

```json
{
    "elements": {
        "div": ["id", "class"],
        "b": [],
        "i": []
    }
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

Create JSON configuration. Below are the currently supported options

| key | value type | default | description |
|-----|------------|---------|-------------|
| stripComments | `boolean` | `false` | Whether or not to strip comment nodes |
| stripWhitespace | `boolean` | `false` | Whether or not to strip whitespace (leading and trailing tabs or spaces) |
| elements| `Object` | `{}` | a list of K-V pairs where the keys are whitelisted element tags and the values are arrays of whitelisted attribtues for that element |

```json
{
    "stripComments": true,
    "stripWhitespace": true,
    "elements": {
        "html": ["xmlns"],
        "head": [],
        "body": [],
        "div": ["id", "class"],
    }
}
```

Create a `sanitize.Whitelist` object from a json file with `sanitize.WhitelistFromFile(filepath string)` or from a []byte with `sanitize.NewWhitelist(byteArray []byte)` and use it to parse some HTML:

```go
whitelist, err := sanitize.WhitelistFromFile("./path/to/file.json")
// or create from a json []byte
// whitelist, err := sanitize.NewWhitelist(byteArray)

f, _ := os.Open("./path/to/example.html")
sanitized, _ := whitelist.SanitizeRemove(f) // takes any io.Reader

fmt.Printf("sanitized html: %d", sanitized)
```

### Supported operations

```go
whitelist, err := sanitize.WhitelistFromFile("./path/to/file.json")
f, _ := os.Open("./path/to/example.html")

// sanitize a full HTML document by removing
// non-whitelisted elements and attributes
sanitized, _ := whitelist.SanitizeRemove(f)

// sanitize a full HTML document by reattaching
// the children of non-whitelisted elements to the
// non-whitelisted parent; also removes non whitelisted
// attributes for any element
sanitized, _ := whitelist.SanitizeUnwrap(f)

// sanitize an HTML document fragment (ie no html,
// head, or body tags) by removing
// non-whitelisted elements and attributes
sanitized, _ := whitelist.SanitizeRemoveFragment(f)

// sanitize an HTML document fragment (ie no html,
// head, or body tags) by reattaching
// the children of non-whitelisted elements to the
// non-whitelisted parent; also removes non whitelisted
// attributes for any element
sanitized, _ := whitelist.SanitizeUnwrapFragment(f)

```

### Steps to 1.0
- [x] Support sanitization that unwraps non-whitelisted nodes, allowing the text and/or whitelisted subtree through
- [x] Whitelist-level configuration options (eg. `stripWhitespace`)
- [x] Efficient attribute checking by not allocating a new slice on every whitelisted attribute for an element
- [x] Support sanitization of HTML fragments (instead of just full documents)
- [ ] Support non `string` type attribute values
- [x] Refactor configuration parsing to have []byte interface instead of expecting a filepath
- [ ] Create sane defaults

### Known Issues

### Contributing

Head over to the [issues page](https://github.com/maxwells/go-html-sanitizer/issues) or open a [pull request](https://github.com/maxwells/go-html-sanitizer/pulls). Please ensure your code is documented, all existing tests pass, and any new features have tests before submitting a pull request. If you want to check in whether a pull request for a new feature would be accepted, feel free to open an issue.

### License

MIT