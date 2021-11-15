## How to Contribute
At NS1, it makes us very happy to get pull requests, and we appreciate the time
and effort it takes to put one together. Below are a few guidelines that can
help get your pull request approved and merged quickly and easily.

It's important to note that the guidelines are recommendations, not
requirements. You may submit a pull request that does not adhere to the
guidelines, but ideation and development may take longer or be more involved.

* Avoid getting nitpicked for basic formatting. Match the style of any given
  file as best you can - tabs or spaces as appropriate, curly bracket
  placement, indentation style, line length, etc. If there are any linters
  mentioned in docs or Makefile, please run them.

* Avoid changes that have any possibility of breaking existing uses. Some of
  these codebases have many of users, doing creative things with them. Breaking
  changes can cause significant challenges for other users. It's important to
  note that some projects are dependencies of other projects, and changes may
  require cross-code base coordination. It can be challenging to identify if a
  small change will have large ramifications, so we mainly ask that you keep
  this in mind when writing or modifying code. Do your best to preserve
  interfaces, and understand that NS1 may need to reject pull requests if they
  would jeopardize platform stability or place an undue burden on other users.

* If there are unit and/or integration tests, keeping the test suites passing
  is a must before we can merge. And nice tests around the changes will
  definitely help get a patch merged quickly.

* Ensure that any documentation that is part of the project is updated as part
  of the pull request. This helps to expedite the merge process.

* Be considerate when introducing new dependencies to a project. If a new
  dependency is necessary, try to choose a well-known project, and to pin the
  version as specifically as possible.

## Opening an issue

Below are a few guidelines for opening an issue on one of NS1's open-source
repositories. Adhering to these guidelines to the best of your ability will
help ensure that requests are resolved quickly and efficiently.

* **Be specific about the problem.** When describing your issue, it's helpful
  to include as many details as you can about the expected behavior versus
  what happened or is happening instead.
* **Be specific about your objective.** Help us understand exactly what you are
  trying to accomplish so that our developers have a clear understanding about
  the particular problem you are trying to solve. In other words, help us
  avoid the "[XY Problem](http://xyproblem.info)".
* **Indicate which products, software versions, and programming languages you
  are using.** In your request, indicate which NS1 product(s) you're using and,
  if relevant, which versions you are running. Also include any third-party
  software and versions that are relevant to the issue. If not obvious, include
  which programming language(s) you are using.
* **If possible, provide a reproducible example of the problem.** This allows
  us to better examine the issue and test solutions.
* **If possible, include stack/error traces.** Note: ensure there is no
  sensitive included in your stack/error traces.

## Closing an issue

* If an issue is closed by a commit, reference the relevant PR or commit when
  closing.
* If an issue is closed by NS1 for any reason, you should expect us to include
  a reason for it.
* If the fix does not work or is incomplete, you are welcome to re-open or
  recreate the issue. When doing so, it's important to be clear about why the
  previous fix was inadequate, to clarify the previous problem statement,
  and/or to modify the scope of the request. Please keep in mind that project
  status consideration or conflicting priorities may require us to close or
  defer work on the new or reopened issue. If that happens, feel free to reach
  out of support@ns1.com for more information.

## Tags on issues

In some projects we (NS1) may apply basic tags on some issues, for
organizational purposes. Note: we do not use tags to indicate timelines or
priorities.

Here are definitions for the most common tags we use:

* **BUG** - This tag confirms that the issue is a bug we intend to fix. The
  issue will remain open until it is resolved.
* **ENHANCEMENT** - This tag indicates that we have categorized the issue as a
  feature request. Depending on priorities and timelines, we may close issues
  with this tag and track them in our internal backlog instead.
