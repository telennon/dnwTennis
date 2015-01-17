<!---
<div class="panel panel-default">
	<div class="panel-heading">
    <h3 class="panel-title">My Heading</h3>
    <div class="panel-body">
    	<h1>Some Stuff</h1>
    </div>
    <p><a href='#/view2' class='btn btn-primary btn-lg' role='button'>Sign Up &raquo;</a></p>
</div>
-->

<?php
      $host = "mysql.lennonkeegan.com";
      $dbName = 'dnwtennis';
      $resTBL = "payors";
      $attTBL = "campers";
      $supTBL = "campSetup";
      $user = 'dnw0001';
      $pass = 'IslandTime';



      try {
        //echo "<p>Making connection to " . $dbName . "</p>";
        $DBH = new PDO("mysql:host=$host;dbname=$dbName", $user, $pass);
        //echo "<p>Setting up exception attribute</p>";
        $DBH->setAttribute( PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION );    
        //echo "<p>Making connection to " . $dbName . "</p>";
        $DBH = new PDO("mysql:host=$host;dbname=$dbName", $user, $pass);
        //echo "<p>Setting up exception attribute</p>";
        $DBH->setAttribute( PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION );
        $STH = $DBH->prepare("SELECT * FROM $supTBL WHERE campyear='2014'");
        if ($STH->execute()) {
          while ($row = $STH->fetch(PDO::FETCH_ASSOC)) {
            $cmp_Year = $row["campyear"];
            $cmp_Cost = $row["campcost"];
            $cmp_C1Dates = $row["camp1dates"];
            $cmp_C2Dates = $row["camp2dates"];
            $cmp_C3Dates = $row["camp3dates"];
            $cmp_SignupStart = new DateTime($row["signupstart"]);
            $cmp_SignupEnd = new DateTime($row["signupend"]);
            $cmp_NoRefunds = $row["norefundsafter"];
          }
        }
      }
      catch(PDOException $e) {
        $dbError = $e->getMessage();
        echo $dbError;
      }
  ?>
    <!-- Fixed navbar
    <div class="navbar navbar-inverse navbar-fixed-top" role="navigation">
      <div class="container">
        <div class="navbar-header">
          <button type="button" class="navbar-toggle" data-toggle="collapse" data-target=".navbar-collapse">
            <span class="sr-only">Toggle navigation</span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
          </button>
          <a class="navbar-brand" href="/~tom/dnwtennis-v2/">DNW Tennis</a>
        </div>
        <div class="navbar-collapse collapse">
          <ul class="nav navbar-nav">
            <li><a href="#Email">Email Coordinator</a></li>
            <li><a href="signup.php">Enter Lottery</a></li>            
          </ul>
        </div>
      </div>
    </div>
    -->

    <div class="container theme-showcase" role="main">

      <!-- Main jumbotron for a primary marketing message or call to action
      <div class="jumbotron">
        <h1>2014 DNW Tennis Camp</h1>
        <h3>Online Sign Up Form</h3>   -->
      
      <div class="panel panel-default">
        <div class="panel-body">
          <div class="col-md-5">
            <dl class="dl-horizontal">
                <dt>2014 Camp Dates</dt>
                <dd>&nbsp</dd>
                <dt>Camp 1</dt>
                <dd><?php echo "$cmp_C1Dates"; ?></dd>
                <dt>Camp 2</dt>
                <dd><?php echo "$cmp_C2Dates"; ?></dd>
                <dt>Camp 3</dt>
                <dd><?php echo "$cmp_C3Dates"; ?></dd>
            </dl>
          </div>

          <div class="col-md-3">
            <dl>
              <dt>2014 Camp Costs</dt>
                <dd><?php echo "$cmp_Cost"; ?> per camper</dd>
              <dt>&nbsp</dt>
              <dt>&nbsp</dt>
            </dl>
          </div>

          <div class="col-md-4">
           <dl>
              <dt>Sign Up Deadline</dt>
                <dd><?php echo $cmp_SignupEnd->format('m-d-Y'); ?></dd>
              <dt>&nbsp</dt>
              <dt>&nbsp</dt>
            </dl>
          </div>
        </div>  <!-- Panel Body -->
      </div>    <!-- Panel -->

      &nbsp
      <?php
        $cur_Time = new DateTime("now");
        if ($cur_Time > $cmp_SignupEnd) {
          echo "<div class=\"alert alert-warning\">";
            echo "<strong>The DNW Tennis Camp lottery has ended.</strong><br />Camps are now being filled on a space available basis. Please complete the signup form. Trish will contact you to advise whether space is available.";
          echo "</div>";
        }

        if ($cur_Time < $cmp_SignupStart) {
          $strt = $cmp_SignupStart->format('m-d-Y');
          echo "<div class=\"alert alert-warning\">";
            echo "<strong>The DNW Tennis Camp lottery opens " . $strt . ".</strong><br /> Please come back then to sign up for the lottery";
          echo "</div>";
        }
      ?>

      <p class="lead">Welcome to the 2014 Decatur Northwest Tennis Camp lottery sign up page. Evan Hundley will be returning this year to offer three tennis camps. Should demand exceed available spots in the camps, DNW family members will be given priority over guests. Each camp is limited to 32 people. The <?php echo $cmp_Cost; ?> camp fee will be billed to your DNW account on June 1st . <strong>No refunds for cancellations will be made after <?php echo $cmp_NoRefunds ?>.</strong>
      </p>

      <p class="lead"> A lottery system will be used if more homeowners sign up for a particular camp than places available. <strong>Sign-ups received after <?php echo $cmp_SignupEnd->format('m-d-Y');  ?> will not be included in the initial drawing.</strong> Please make sure to indicate your first camp choice and a second choice. Entire families will be signed up before moving on to the next form. Following homeowners, guests will be added in the order that forms were drawn. A wait list for each camp will be created as needed. You will be notified of the results of the drawing on or around May 1st. If you have any difficulties with the online submission or have any questions, feel free to contact me by email or phone. 
      </p>
  <div class="row">
    <div class="col-md-8">
        <address>
          <strong>Trish Keegan</strong><br>
          &nbsp&nbsp&nbsp<a href="mailto:Trish@DNWTennis.com">Trish@DNWTennis.com</a><br>
          &nbsp&nbsp&nbsp206.406.4709
        </address>
    </div>
    <?php
      if ($cur_Time >= $cmp_SignupStart) {
        echo("<div class='col-md-4'>");
          echo("<br><p><a href='#/view2' class='btn btn-primary btn-lg' role='button'>Sign Up &raquo;</a></p>");
        echo("</div>");
      }
    ?>

  </div>


<!--      </div>  -->

    </div> <!-- /container -->