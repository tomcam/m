{{ toc }}

## Special effects

### Last item in a header unordered list shows outlined, simulating a button

```
/*                                                                      
 * --------------------------------------------------
 * Special feature: Last item in list is styled 
 * to look like a button. 
 * See also theme-light.css and theme-dark.css
 * --------------------------------------------------
 */

header > ul > li:last-child > a {border: 2px solid var(--header-fg);color:var(--header-fg);padding:.5rem;}
```

## Sidebar 

Show unordered list in a sidebar as circled numbers. Works well for a news-type section
```
/*
 * --------------------------------------------------
 * Sidebar unordered list for breaking news
 * shows:
 * - Item preceded by number in circle
 * - Item is bold
 * - Indented item in normal text with bottom padding
 *
 * Example usage (note 2nd level of indentation):
 *
 *  * Item 1
 *     - More about item 1
 *   * Item 2
 *     - More about item 2
 *   * Item 3
 *     - More about item 3
 *
 * --------------------------------------------------
 */

aside > ul {
  counter-reset:li;
  list-style-type:none;
  font-size:1rem;
  line-height:1.2rem;
  padding-left:1em;
  border-top:none;
}

aside > ul li {
  list-style-type:none;
  border:none;
  font-weight:normal;
}

aside > ul > li > ul > li {
  line-height:1.5em;
  padding-bottom:1em;
}
aside > ul >li {
  font-weight:bold;
  position:relative;
  padding: .5em 1em 0rem 3em;
  border:none;
} 
aside > ul > li:before {
  background-color:var(--fg);
  content:counter(li;
  counter-increment:li;
  height:2rem;
  line-height:2rem;
  width:2em;
  border-radius:50%;
  color: var(--bg);
  text-align:center;
  position:absolute;
  left:0;
}
```

### Create table of contents in sidebar that shows as boxes of text

```
/*
 * --------------------------------------------------
 * Sidebar unordered list shows as boxes, without
 * indentation--it's for table of contents
 * --------------------------------------------------
 */
aside > ul {
  background-color:whitesmoke;
  margin-right:1em;
  border-collapse:collapse;
}

aside > ul li {
  list-style-type:none;
  margin-left:0;
  /* Border bottom stretches across column at all levels */
  border-bottom: 1px solid gray;
}

aside > ul li a {
  padding-left:.5em;padding-right:.5em;
  text-decoration:none;
  line-height:1em;
}

aside > ul li a:active, aside > ul li a:hover {
  color:blue;
}
```


