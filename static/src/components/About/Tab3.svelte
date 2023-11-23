<script>
</script>

<attributes>
  <div>
    <p>
      Additional attributes can be passed when requesting a ticket. These
      attributes will be associated with the ticked and the user, for which the
      ticket is generated.
    </p>
    <p>
      The attributes are not permanent and will be valid only for the current
      session.
    </p>
    <p>
      As mentioned in the attributes text box the attributes should be provided
      in the form of a JSON array. Each property name will be the attribute name
      and each value will be the attribute value.
    </p>
    <pre>
  {`[
    { "group": "some group" },
    { "group": "another group" },
    { "otherProperty": "some value" },
  ]`}
</pre>
    <div>
      In the example above three attributes will be set:
      <ul>
        <li><code>group</code> with value <code>some group</code></li>
        <li><code>group</code> with value <code>another group</code></li>
        <li><code>otherProperty</code> with value <code>some value</code></li>
      </ul>
    </div>
  </div>
</attributes>

<style>
  div {
    text-align: left;
  }

  p {
    font-size: 17px;
  }
</style>
