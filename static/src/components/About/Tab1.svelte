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
</div>

<style>
  div {
    text-align: left;
  }

  p {
    font-size: 17px;
  }
</style>
