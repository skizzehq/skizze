/*
Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

/*
Package cml implments the Count-Min-Log datastructure.

It is based on Count-Min-Log sketch: Approximately counting with approximate counters - Guillaume Pitel & Geoffroy Fouquier

http://iswag-symposium.org/2015/pdfs/shortpaper1.pdf

TL;DR: Count-Min-Log Sketch for improved Average Relative Error on low frequency events

Count-Min Sketch is a widely adopted algorithm for approximate event counting in large scale processing. However, the original version of the Count-Min-Sketch (CMS) suffers of some deficiences, especially if one is interested in the low-frequency items, such as in text- mining related tasks. Several variants of CMS have been proposed to compensate for the high relative error for low-frequency events, but the proposed solutions tend to correct the errors instead of preventing them. In this paper, we propose the Count-Min-Log sketch, which uses logarithm-based, approximate counters instead of linear counters to improve the average relative error of CMS at constant memory footprint.
*/
package cml
