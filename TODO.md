# todo

We should use the parser feature of `golang.org/x/net/html`.
It might not always parse correctly the HTML when it is malformed and might skip some nodes.
But it should still do the job compared to goquery that is based of this package.
Of what I see, goquery brings features like a `Find` function that might not be needed here.

We need to extract the body.
We need to be able to ignore some tags.
We should be able to recreate a simplified version of the HTML.
