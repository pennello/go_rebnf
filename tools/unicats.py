# -*_ coding: utf-8 -*-

# chris 073015

import sys

from HTMLParser import HTMLParser
from os.path import commonprefix
from contextlib import closing
from urllib2 import urlopen

base = 'http://www.fileformat.info/info/unicode/category'

htmlunescape = HTMLParser().unescape

def getchars(cat):
  url = '%s/%s/list.htm' % (base,cat)
  chars = []
  with closing(urlopen(url)) as page:
    for line in page:
      if 'U+' not in line: continue
      idx = line.index('U+')
      idx += 2
      idx2 = line.index('<',idx)
      point = int(line[idx:idx2],16)
      descr = page.next().strip()
      # Eliminate <td> and </td>
      descr = descr[4:-5]
      descr = htmlunescape(descr)
      chars.append((point,descr))
  return chars

def main(self,*cats):
  chars = []
  for cat in cats:
    chars.extend(getchars(cat))

  # Consolidate contiguous ranges.
  chars.sort()
  ranges = []
  last = start = None
  descrs = []
  for point,descr in chars:
    if start is None:
      start = point
    elif point != last + 1:
      end = last
      descr2 = commonprefix(descrs)
      if not descr2: descr2 = '?'
      ranges.append(((start,end),descr2))
      start = point
      descrs = []

    last = point
    descrs.append(descr)

  for (start,end),descr in ranges:
    out = ''
    if start == end:
      out += r'"\U%08x"' % start
      out += ' ' * 15
    else:
      out += r'"\U%08x" … "\U%08x"' % (start,end)
    out += ' | // %s' % descr
    print out

main(*sys.argv)
