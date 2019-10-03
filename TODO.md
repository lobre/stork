# todo

We should use the parser feature of `golang.org/x/net/html`.
It might not always parse correctly the HTML when it is malformed and might skip some nodes.
But it should still do the job compared to goquery that is based of this package.
Of what I see, goquery brings features like a `Find` function that might not be needed here.

We need to extract the body.
We need to be able to ignore some tags.
We should be able to recreate a simplified version of the HTML.

---

Formatting HTML is a pretty complicated task. We will try to implement a simple version of taste here.
We first need to determine the type of HTML element to know if we want to keep it / indent it or not.
A good reference to determine this is this document: https://developer.mozilla.org/en-US/docs/Web/Guide/HTML/Content_categories
Or for a simpler document (prior HTML5): https://www.w3schools.com/html/html_blocks.asp

We can:
 - remove metadata elements
 - keep phrasing / inline elements compact (strip whitespaces and don't add newlines before / after)
 - add newlines for block elements
 - make sure <pre> does not contain added spaces

The strategy will apply on text elements (because they will contain the newlines). That means we will have to check the parent / previous / next element tags.

---

Another directly for spacing.
Render() will render HTML with good spacing according to HTML rules.
Try to study Render() to be sure about this assumption.

So what is important is to clean the Text elements of the Node tree in order to simplify the Text and Markdown outputs.

See https://github.com/jaytaylor/html2text/blob/master///html2text.go#L392:36

---

Should we use context in recursive functions?
