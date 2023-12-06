<div>
  <p>
    The basic idea is to have a list of test users in Qlik Sense and the service
    can generate web tickets for them. Once the ticket is available the
    developer/tester can use the ticket to login to Qlik Sense as the test user
    and "see" throw the eyes of that user.
  </p>
  <p>
    <strong>
      This service is not meant to be used in production environment. Its main
      purpose is to help the developers/testers to verify their work in the
      lower environments (PRE-PROD, UAT, DEV, TEST etc.)
    </strong>
  </p>

  <p>The generated test users:</p>
  <ul>
    <li>
      can be added to section access with specific permissions (which can be
      altered)
    </li>
    <li>
      test the security rules behavior without constantly asking the real users
      to "refresh again?" in order to see the security rule implication
    </li>
  </ul>

  <div>
    Please refer to the complete <a
      target="_blank"
      href="https://github.com/Informatiqal/qlik-test-users-tickets#readme"
      >documentation</a
    > for more information
  </div>
</div>

<style>
  div {
    text-align: left;
  }

  p {
    font-size: 17px;
  }

  a[target="_blank"]::after {
    content: "";
    width: 1em;
    height: 1em;
    margin: 0 0.05em 0 0.1em;
    background: url(data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAxNiAxNiIgd2lkdGg9IjE2IiBoZWlnaHQ9IjE2IiBmaWxsPSIjNjQ2Y2ZmIj48cGF0aCBkPSJNOSAyTDkgMyAxMi4zIDMgNiA5LjMgNi43IDEwIDEzIDMuNyAxMyA3IDE0IDcgMTQgMlpNNCA0QzIuOSA0IDIgNC45IDIgNkwyIDEyQzIgMTMuMSAyLjkgMTQgNCAxNEwxMCAxNEMxMS4xIDE0IDEyIDEzLjEgMTIgMTJMMTIgNyAxMSA4IDExIDEyQzExIDEyLjYgMTAuNiAxMyAxMCAxM0w0IDEzQzMuNCAxMyAzIDEyLjYgMyAxMkwzIDZDMyA1LjQgMy40IDUgNCA1TDggNSA5IDRaIi8+PC9zdmc+)
      no-repeat;
    background-size: contain;
    display: inline-block;
    vertical-align: text-bottom;
    padding-bottom: 2px;
  }
</style>
