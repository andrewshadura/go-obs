OBS API client in Go
====================

go-obs is an API client for OBS, [Open Build Service][].

go-obs currently only covers a tiny part of the OBS API. It also provides a
very basic command-line client tool.

Please note that this project has nothing to do with another OBS which is
short for [Open Broadcaster Software][]. If you’re looking for a Go client
for that OBS, have a look at [obsws][].

API coverage
------------

The client currently only supports most of the user and group manipulation operations.

License
-------

Licensed under the Apache License, Version 2.0 (the ‘License’);
you may not use this software except in compliance with the License.
You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>;
it is also included as [LICENSE](./LICENSE) with this software.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

[Open Build Service]: https://openbuildservice.org/
[Open Broadcaster Software]: https://obsproject.com/
[obsws]: https://github.com/christopher-dG/go-obs-websocket
