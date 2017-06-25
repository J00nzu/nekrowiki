# Hello

What the fuck

ハオ・ドス・ジュニコード・ワーク 

- [x] @mentions, #refs, [links](index), **formatting**, and <del>tags</del> supported
- [x] list syntax required (any unordered or ordered list supported)
- [x] this is a complete item
- [ ] this is an incomplete item

First Header | Second Header
------ | ---
Content from cell 1 | Content from cell 2
Content in the first column | Content in the second column

```C++
#include <iostream>
#include <vector>
#include <stdexcept>

int main() {
    try {
        std::vector<int> vec{3,4,3,1};
        int i{vec.at(4)}; // Throws an exception, std::out_of_range (indexing for vec is from 0-3 not 1-4)
    }

    // An exception handler, catches std::out_of_range, which is thrown by vec.at(4)
    catch (std::out_of_range& e) {
        std::cerr << "Accessing a non-existent element: " << e.what() << '\n';
    }

    // To catch any other standard library exceptions (they derive from std::exception)
    catch (std::exception& e) {
        std::cerr << "Exception thrown: " << e.what() << '\n';
    }

    // Catch any unrecognised exceptions (i.e. those which don't derive from std::exception)
    catch (...) {
        std::cerr << "Some fatal error\n";
    }
}
```


:sparkles: :camel: :boom:

# This is an \<h1> tag
## This is an \<h2> tag
###### This is an \<h6> tag

*This text will be italic*
_This will also be italic_

**This text will be bold**
__This will also be bold__

_You **can** combine them_

This_Should_Work_Too

![TEST IMAGE](/uploads/bsglogo3.png)